package agentproxy

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gliderlabs/ssh"
	"github.com/inconshreveable/muxado"
	gossh "golang.org/x/crypto/ssh"
)

var TmpDir = "/tmp/sock"
var Image = "gliderlabs/ssh-agent-proxy"

type AgentProxy struct {
	SocketPath  string
	ContainerID string

	docker client.APIClient
	sess   ssh.Session
	stream *types.HijackedResponse
}

func NewAgentProxy(docker client.APIClient, sess ssh.Session) (*AgentProxy, error) {
	sessID := sess.Context().Value(ssh.ContextKeySessionID).(string)[:12]
	socketPath := path.Join(TmpDir, fmt.Sprintf("auth-agent.%s", sessID), "listener.sock")
	res, err := docker.ContainerCreate(context.Background(), &container.Config{
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
