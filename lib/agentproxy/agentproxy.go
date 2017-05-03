package agentproxy

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/ssh"
	"github.com/inconshreveable/muxado"
	gossh "golang.org/x/crypto/ssh"
)

func init() {
	com.Register("agentproxy", &AgentProxy{},
		com.Option("tmpdir", "/tmp/sock", "auth-agent socket directory"),
		com.Option("image", "gliderlabs/ssh-agent-proxy", "Docker image to use for the agent proxy"),
	)
}

type AgentProxy struct {
	SocketPath  string
	ContainerID string

	docker client.APIClient
	sess   ssh.Session
	stream *types.HijackedResponse
}

// ImageExists and InspectImage taken from libcompose
// Exists return whether or not the service image already exists
func Exists(ctx context.Context, clt client.ImageAPIClient, image string) (bool, error) {
	_, err := InspectImage(ctx, clt, image)
	if err != nil {
		if client.IsErrImageNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// InspectImage inspect the specified image (can be a name, an id or a digest)
// with the specified client.
func InspectImage(ctx context.Context, client client.ImageAPIClient, image string) (types.ImageInspect, error) {
	imageInspect, _, err := client.ImageInspectWithRaw(ctx, image)
	return imageInspect, err
}

func NewAgentProxy(docker client.APIClient, sess ssh.Session) (*AgentProxy, error) {
	var (
		TmpDir = com.GetString("tmpdir")
		Image  = com.GetString("image")
	)
	sessID := sess.Context().Value(ssh.ContextKeySessionID).(string)[:12]
	socketPath := path.Join(TmpDir, fmt.Sprintf("auth-agent.%s", sessID), "listener.sock")
	ctx := context.Background()

	exists, err := Exists(ctx, docker, Image)
	if err != nil {
		return nil, err
	}
	if !exists {
		res, err := docker.ImagePull(ctx, Image, types.ImagePullOptions{})
		if err != nil {
			return nil, err
		}
		io.Copy(ioutil.Discard, res)
		res.Close()
	}

	res, err := docker.ContainerCreate(ctx, &container.Config{
		Image:        Image,
		Cmd:          []string{socketPath},
		OpenStdin:    true,
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		StdinOnce:    true,
	}, nil, nil, "")
	if err != nil {
		return nil, err
	}
	return &AgentProxy{
		SocketPath:  socketPath,
		ContainerID: res.ID,
		docker:      docker,
		sess:        sess,
	}, nil
}

func (ap *AgentProxy) Shutdown() error {
	if ap.ContainerID != "" {
		ctx := context.Background()
		err := ap.docker.ContainerRemove(ctx, ap.ContainerID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			return err
		}
	}
	if ap.stream != nil {
		ap.stream.Conn.Close()
	}
	return nil
}

func (ap *AgentProxy) Start() error {
	ctx := context.Background()
	stream, err := ap.docker.ContainerAttach(ctx, ap.ContainerID, types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: true,
	})
	if err != nil {
		return err
	}
	ap.stream = &stream
	pr, pw := io.Pipe()
	go func() {
		stdcopy.StdCopy(pw, os.Stderr, stream.Reader)
		pw.Close()
	}()
	tunnel := muxado.Server(struct {
		io.Reader
		io.WriteCloser
	}{pr, stream.Conn}, nil)
	go func() {
		for {
			stream, err := tunnel.AcceptStream()
			if err != nil {
				break
			}
			go ap.proxyStream(stream)
		}
	}()
	return ap.docker.ContainerStart(ctx, ap.ContainerID, types.ContainerStartOptions{})
}

func (ap *AgentProxy) proxyStream(stream muxado.Stream) {
	defer stream.Close()
	sshConn := ap.sess.Context().Value(ssh.ContextKeyConn).(gossh.Conn)
	channel, reqs, err := sshConn.OpenChannel("auth-agent@openssh.com", nil)
	if err != nil {
		panic(err)
	}
	defer channel.Close()
	go gossh.DiscardRequests(reqs)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		io.Copy(stream, channel)
		stream.CloseWrite()
		wg.Done()
	}()
	go func() {
		io.Copy(channel, stream)
		channel.CloseWrite()
		wg.Done()
	}()
	wg.Wait()
}
