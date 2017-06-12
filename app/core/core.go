package core

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-units"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/gliderlabs/ssh"
	"github.com/pkg/errors"

	"github.com/gliderlabs/cmd/app/billing"
	"github.com/gliderlabs/cmd/lib/agentproxy"
	"github.com/gliderlabs/cmd/lib/crypto"
	"github.com/gliderlabs/cmd/lib/docker"
	"github.com/gliderlabs/cmd/lib/release"
)

const ServerSoftware = "cmd.io"

// Token used to provide access to non-github users
type Token struct {
	Key         string
	Description string
	User        string
	LastUsedIP  string
	LastUsedOn  time.Time
}

func (t *Token) Validate() error {
	if t.Key == "" {
		return fmt.Errorf("token Key required")
	}

	if t.User == "" {
		return fmt.Errorf("token User required")
	}
	return nil
}

// Command is a the definition for a runnable command on cmd
type Command struct {
	Name        string
	User        string
	Source      string
	Environment map[string]string `dynamodbav:",omitempty"`
	ACL         []string          `dynamodbav:",stringset,omitempty"`
	Admins      []string          `dynamodbav:",stringset,omitempty"`
	Description string            `dynamodbav:",omitempty"`

	Changed bool `dynamodbav:"-"`

	docker client.APIClient
}

// Docker will return a configured docker client
func (c *Command) Docker() client.APIClient {
	if c.docker == nil {
		c.docker = docker.Client()
	}
	return c.docker
}

// SetEnv for command
func (c *Command) SetEnv(key, val string) {
	if c.Environment == nil {
		c.Environment = make(map[string]string)
	}
	c.Environment[key] = val
}

// Env returns config in a `k=v` format without any cmd specific keys
func (c *Command) Env() (env []string) {
	env = append(env, []string{
		fmt.Sprintf("SERVER_SOFTWARE=%s", ServerSoftware),
		fmt.Sprintf("CMD_CHANNEL=%s", release.Channel()),
		fmt.Sprintf("CMD_VERSION=%s", release.DisplayVersion()),
	}...)
	for k, v := range c.Environment {
		if strings.HasPrefix(k, "io.cmd") {
			continue
		}
		env = append(env, fmt.Sprintf("%s=%s", k, crypto.Decrypt(v)))
	}
	return
}

func (c *Command) HasAccess(user string) bool {
	if c.User == user {
		return true
	}
	for _, u := range c.ACL {
		if u == user {
			return true
		}
	}
	for _, u := range c.Admins {
		if u == user {
			return true
		}
	}
	return false
}

func (c *Command) IsAdmin(user string) bool {
	if c.User == user {
		return true
	}
	for _, u := range c.Admins {
		if u == user {
			return true
		}
	}
	return false
}

func (c *Command) image() string {
	return fmt.Sprintf("%s-%s", c.User, c.Name)
}

func parseSource(src []byte) (img string, pkgs []string, body []byte, err error) {
	if !bytes.HasPrefix(src, []byte("#!cmd")) {
		err = errors.Errorf("invalid source: first line must start with `#!cmd`")
		return
	}
	adv, tok, _ := bufio.ScanLines(src, false)
	if tok == nil {
		err = errors.Errorf("unable to parse input")
		return
	}
	body = src[adv:]
	parts := strings.Fields(string(tok))
	if len(parts) > 1 {
		img = parts[1]
	}
	if len(parts) > 2 {
		pkgs = parts[2:]
	}
	return
}

func getBuildCtx(img string, pkgs []string, body []byte) (ctx map[string][]byte, err error) {
	ctx = map[string][]byte{}
	if img != "alpine" {
		err = errors.Errorf("unsupported image: currently alpine is the only supported image")
		return
	}

	var dockerfile bytes.Buffer
	fmt.Fprintln(&dockerfile, "FROM", img)
	if len(pkgs) != 0 {
		fmt.Fprintln(&dockerfile,
			"RUN apk --no-cache add", strings.Join(pkgs, " "))
	}
	adv, entrypoint, _ := bufio.ScanLines(body, false)
	if entrypoint == nil {
		err = errors.Errorf("unable to parse body")
		return
	}
	entrypoint = bytes.TrimPrefix(entrypoint, []byte("#!"))
	if len(body)-adv != 0 {
		ctx["entrypoint"] = body
		fmt.Fprintln(&dockerfile, "COPY ./entrypoint ./bin/entrypoint")
		entrypoint = []byte("/bin/entrypoint")
	}
	fmt.Fprintln(&dockerfile, "WORKDIR", "/cmd")
	fmt.Fprintln(&dockerfile, "ENTRYPOINT", `["`+string(entrypoint)+`"]`)
	ctx["Dockerfile"] = dockerfile.Bytes()
	return
}

func (c *Command) Build() error {
	img, pkgs, body, err := parseSource([]byte(c.Source))
	if err != nil {
		return err
	}

	buildCtx, err := getBuildCtx(img, pkgs, body)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	for name, b := range buildCtx {
		err = tw.WriteHeader(&tar.Header{
			Name: name,
			Mode: 0700,
			Size: int64(len(b)),
		})
		if err != nil {
			return err
		}
		if _, err = tw.Write(b); err != nil {
			return err
		}
	}
	if err = tw.Close(); err != nil {
		return err
	}

	ctx := context.Background()
	r := bytes.NewReader(buf.Bytes())
	resp, err := c.Docker().ImageBuild(ctx, r, types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{c.image()},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// read and discard all output to ensure we don't return before the image
	// is built and tagged
	io.Copy(ioutil.Discard, resp.Body)
	return nil
}

// Pull and tag image for command
func (c *Command) Pull(ctx context.Context) error {
	res, err := c.Docker().ImagePull(ctx, c.Source, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, res)
	res.Close()

	img, _, err := c.Docker().ImageInspectWithRaw(ctx, c.Source)
	if err != nil {
		return err
	}

	if maxSize := billing.ContextPlan(ctx).ImageSize; img.Size > maxSize {
		c.Docker().ImageRemove(ctx, c.Source, types.ImageRemoveOptions{}) // Do something with error
		return errors.Errorf("image size excedes plan limit of: %s with: %s",
			units.BytesSize(float64(maxSize)),
			units.BytesSize(float64(img.Size)))
	}
	return c.Docker().ImageTag(ctx, c.Source, c.image())
}

// Run a command in a container attaching input/output to ssh session
func (c *Command) Run(sess ssh.Session, args []string) int {
	var err error
	if strings.HasPrefix(c.Source, "#!") {
		err = c.Build()
	} else {
		err = c.Pull(sess.Context())
	}
	if err != nil {
		fmt.Fprintln(sess.Stderr(), err.Error())
		return 255
	}

	status, err := c.run(sess, args)
	if err != nil {
		fmt.Fprintln(sess.Stderr(), err.Error())
		return status
	}

	return status
}

func (c *Command) run(sess ssh.Session, args []string) (int, error) {
	pty, winCh, isPty := sess.Pty()
	client := c.Docker()
	env := append([]string{
		"REMOTE_ADDR=" + sess.RemoteAddr().String(),
		"USER=" + sess.User(),
		"CMD_NAME=" + sess.Command()[0],
	}, c.Env()...)
	env = append(env, sess.Environ()...)
	if isPty {
		env = append([]string{fmt.Sprintf("TERM=%s", pty.Term)}, env...)
	}
	ctx := sess.Context()
	p := billing.ContextPlan(ctx)
	hostConf := &container.HostConfig{
		AutoRemove: true,
		Resources: container.Resources{
			CPUPeriod: p.CPUPeriod,
			CPUQuota:  p.CPUQuota,
			Memory:    p.Memory,
		},
	}
	if ssh.AgentRequested(sess) {
		proxy, err := agentproxy.NewAgentProxy(client, sess)
		if err != nil {
			return 255, err
		}
		if err := proxy.Start(); err != nil {
			return 255, err
		}
		defer proxy.Shutdown()
		env = append(env, fmt.Sprintf("SSH_AUTH_SOCK=%s", proxy.SocketPath))
		hostConf.VolumesFrom = []string{proxy.ContainerID}
	}
	conf := &container.Config{
		Image:        c.image(),
		Env:          env,
		Cmd:          args,
		Tty:          isPty,
		OpenStdin:    true,
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		StdinOnce:    true,
		Volumes:      make(map[string]struct{}),
	}

	if p.DinD {
		// TODO: actual feature, maybe: https://github.com/gliderlabs/cmd/issues/40
		conf.Volumes["/var/run/docker.sock"] = struct{}{}
		hostConf.Binds = []string{"/var/run/docker.sock:/var/run/docker.sock"}
	}
	res, err := client.ContainerCreate(ctx, conf, hostConf, nil, "")
	if err != nil {
		return 255, err
	}
	defer client.ContainerRemove(ctx, res.ID, types.ContainerRemoveOptions{Force: true})
	containerStream, err := client.ContainerAttach(ctx, res.ID,
		types.ContainerAttachOptions{
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			Stream: true,
		})
	if err != nil {
		return 255, err
	}
	defer containerStream.Close()
	receiveStream := make(chan error, 1)
	go func() {
		var err error
		if isPty {
			_, err = io.Copy(sess, containerStream.Reader)
		} else {
			_, err = stdcopy.StdCopy(sess, sess.Stderr(), containerStream.Reader)
		}
		receiveStream <- err
	}()

	go func() {
		defer containerStream.CloseWrite()
		io.Copy(containerStream.Conn, sess)
	}()

	err = client.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})
	if err != nil {
		return 255, err
	}

	if isPty {
		go func() {
			for win := range winCh {
				err := client.ContainerResize(ctx, res.ID, types.ResizeOptions{
					Height: uint(win.Height),
					Width:  uint(win.Width),
				})
				if err != nil {
					log.Info(errors.Wrap(err, "failed to resize pty"))
					break
				}
			}
		}()
	}

	statusChan := make(chan int64, 1)
	go func() {
		s, err := client.ContainerWait(ctx, res.ID)
		if err != nil {
			log.Info(errors.Wrap(err, "container wait failed"))
		}
		statusChan <- s
	}()

	timeout := time.After(p.MaxRuntime)
	select {
	case <-timeout:
		return 255, billing.ErrMaxRuntimeExceded
	case err := <-receiveStream:
		if err != nil {
			return 255, err
		}
	}

	status := <-statusChan
	return int(status), nil
}
