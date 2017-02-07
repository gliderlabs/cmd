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
	"github.com/docker/docker/pkg/stdcopy"
	units "github.com/docker/go-units"
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/pkg/errors"
	"github.com/progrium/cmd/pkg/dune"
)

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

type Stream struct {
	Stdin  io.ReadCloser
	Stdout io.Writer
	Stderr io.Writer
}

// Command is a the definition for a runnable command on cmd
type Command struct {
	Name        string
	User        string
	Environment map[string]string
	ACL         []string
	Admins      []string
	Description string
	Source      string

	Changed bool

	docker *dune.Client
}

// Docker will return a configured docker client
func (c *Command) Docker() *dune.Client {
	if c.docker != nil {
		return c.docker
	}
	var err error
	c.docker, err = dune.NewClient(com.GetString("host"))
	if err != nil {
		log.Info(errors.Wrap(err, "failed to create new dune client"))
	}
	if c.docker != nil {
		return c.docker
	}
	c.docker, err = dune.NewEnvClient()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create docker env client"))
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
	for k, v := range c.Environment {
		if strings.HasPrefix(k, "io.cmd") {
			continue
		}
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return
}

// MakePrivate by setting ACL to an empty string slice
func (c *Command) MakePrivate() {
	c.ACL = []string{}
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

// AddAccess for user to command
func (c *Command) AddAccess(user string) {
	if c.HasAccess(user) {
		return
	}
	c.ACL = append(c.ACL, user)
}

// RemoveAccess from user to command
func (c *Command) RemoveAccess(user string) {
	var i int
	var u string
	var found bool
	for i, u = range c.ACL {
		if u == user {
			found = true
			break
		}
	}
	if found {
		c.ACL = append(c.ACL[:i], c.ACL[i+1:]...)
	}
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

func (c *Command) AddAdmin(user string) {
	if c.IsAdmin(user) {
		return
	}
	if c.HasAccess(user) {
		c.RemoveAccess(user)
	}
	c.Admins = append(c.Admins, user)
}

func (c *Command) RemoveAdmin(user string) {
	var i int
	var u string
	var found bool
	for i, u = range c.Admins {
		if u == user {
			found = true
			break
		}
	}
	if found {
		c.Admins = append(c.Admins[:i], c.Admins[i+1:]...)
	}
}

func (c *Command) image() string {
	return fmt.Sprintf("%s-%s", c.User, c.Name)
}

func parseSource(src []byte) (img string, pkgs []string, body []byte, err error) {
	img = "alpine" // Use alpine as default image
	adv, tok, _ := bufio.ScanLines(src, false)
	if tok == nil {
		err = errors.Errorf("unable to parse input")
		return
	}
	body = src[adv:]
	parts := strings.Fields(string(tok))
	if len(parts) > 2 {
		img = parts[1]
		pkgs = parts[2:]
		return
	}
	return
}

func getBuildCtx(img string, pkgs []string, body []byte) (ctx map[string][]byte, err error) {
	ctx = map[string][]byte{}
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

	if img != "alpine" {
		return errors.Errorf("unsupported image: currently alpine is the only supported image")
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
	_, err = c.Docker().ImageBuild(ctx, r, types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{c.image()},
	})
	return err
}

// Pull and tag image for command
func (c *Command) Pull() error {
	ctx := context.Background()
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

	if maxSize := Plans[DefaultPlan].ImageSize; img.Size > maxSize {
		c.Docker().ImageRemove(ctx, c.Source, types.ImageRemoveOptions{}) // Do something with error
		return errors.Errorf("image excedes plan size limit of: %s with: %s",
			units.HumanSize(float64(maxSize)),
			units.HumanSize(float64(img.Size)))
	}
	return c.Docker().ImageTag(ctx, c.Source, c.image())
}

func (c *Command) UpdateDescription() error {
	res, _, err := c.Docker().ImageInspectWithRaw(context.Background(), c.image())
	if err != nil {
		return err
	}
	desc := res.ContainerConfig.Labels["io.cmd.description"]
	c.Description = strings.Trim(desc, "\n")
	c.Changed = true
	return nil
}

// Run a command in a container attaching input/output to ssh session
func (c *Command) Run(stream *Stream, user string, args []string) int {
	var err error
	if strings.HasPrefix(c.Source, "#!") {
		err = c.Build()
	} else {
		err = c.Pull()
	}
	if err != nil {
		fmt.Fprintln(stream.Stderr, err.Error())
		return 255
	}
	if err = c.UpdateDescription(); err != nil {
		fmt.Fprintln(stream.Stderr, err.Error())
		return 255
	}

	status, err := c.run(stream, user, args)
	if err != nil {
		fmt.Fprintln(stream.Stderr, err.Error())
		return status
	}

	return status
}

// func (c *Command) run(s ssh.Session, args []string) error {
func (c *Command) run(stream *Stream, user string, args []string) (int, error) {
	// outputStream, errorStream, inputStream := stdout, s.Stderr(), s
	client := c.Docker()
	ctx := context.Background()
	env := append(c.Env(), "USER="+user)
	conf := &container.Config{
		Image:        c.image(),
		Env:          env,
		Cmd:          args,
		OpenStdin:    true,
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		StdinOnce:    true,
		Volumes:      make(map[string]struct{}),
	}

	hostConf := &container.HostConfig{
		Resources: container.Resources{
			CPUPeriod: Plans[DefaultPlan].CPUPeriod,
			CPUQuota:  Plans[DefaultPlan].CPUQuota,
			Memory:    Plans[DefaultPlan].Memory,
		},
	}

	if c.User == "progrium" || c.User == "mattaitchison" {
		conf.Volumes["/var/run/docker.sock"] = struct{}{}
		hostConf.Binds = []string{"/var/run/docker.sock:/var/run/docker.sock"}
		hostConf.Privileged = true
		hostConf.UsernsMode = "host"
	}

	res, err := client.ContainerCreate(ctx, conf, hostConf, nil, "")
	if err != nil {
		return 255, err
	}

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
		_, copyErr := stdcopy.StdCopy(stream.Stdout, stream.Stderr, containerStream.Reader)
		receiveStream <- copyErr
	}()

	inputDone := make(chan struct{})
	go func() {
		io.Copy(containerStream.Conn, stream.Stdin)
		if copyErr := containerStream.CloseWrite(); copyErr != nil {
			log.Debug("Couldn't send EOF: %s", copyErr)
		}
		close(inputDone)
	}()

	statusChan := make(chan int64, 1)
	go func() {
		s, _ := client.ContainerWait(ctx, res.ID)
		statusChan <- s
	}()

	err = client.ContainerStart(ctx, res.ID, types.ContainerStartOptions{})
	if err != nil {
		return 255, err
	}

	maxDur := Plans[DefaultPlan].MaxRuntime
	timeout := time.After(maxDur)
	select {
	case <-timeout:
		return 255, ErrMaxRuntimeExceded
	case err := <-receiveStream:
		if err != nil {
			return 255, err
		}
	case <-inputDone:
		select {
		case err := <-receiveStream:
			if err != nil {
				return 255, err
			}
		}
	}

	status := <-statusChan
	client.ContainerRemove(ctx, res.ID, types.ContainerRemoveOptions{})
	return int(status), nil
}
