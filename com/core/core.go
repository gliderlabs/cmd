package core

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gliderlabs/gosper/pkg/com"
	"github.com/gliderlabs/gosper/pkg/log"
	"github.com/gliderlabs/ssh"
	"github.com/pkg/errors"
	"github.com/progrium/cmd/pkg/dune"
)

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

// Pull and tag image for command
func (c *Command) Pull() error {
	ctx := context.Background()
	res, err := c.Docker().ImagePull(ctx, c.Source, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, res)
	res.Close()
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
func (c *Command) Run(s ssh.Session, args []string) {
	if err := c.Pull(); err != nil {
		fmt.Fprintln(s.Stderr(), err.Error())
		s.Exit(255)
		return
	}
	if err := c.UpdateDescription(); err != nil {
		fmt.Fprintln(s.Stderr(), err.Error())
		s.Exit(255)
		return
	}
	if err := c.run(s, args); err != nil {
		fmt.Fprintln(s.Stderr(), err.Error())
		s.Exit(255)
		return
	}
}

func (c *Command) run(s ssh.Session, args []string) error {
	outputStream, errorStream, inputStream := s, s.Stderr(), s
	client := c.Docker()
	ctx := context.Background()
	conf := &container.Config{
		Image:        c.image(),
		Env:          c.Env(),
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
		return err
	}

	containerStream, err := client.ContainerAttach(ctx, res.ID,
		types.ContainerAttachOptions{
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			Stream: true,
		})
	if err != nil {
		return err
	}
	defer containerStream.Close()

	receiveStream := make(chan error, 1)
	go func() {
		_, copyErr := stdcopy.StdCopy(outputStream, errorStream, containerStream.Reader)
		receiveStream <- copyErr
	}()

	inputDone := make(chan struct{})
	go func() {
		io.Copy(containerStream.Conn, inputStream)
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
		return err
	}

	maxDur := Plans[DefaultPlan].MaxRuntime
	timeout := time.After(maxDur)
	select {
	case <-timeout:
		return ErrMaxRuntimeExceded
	case err := <-receiveStream:
		if err != nil {
			return err
		}
	case <-inputDone:
		select {
		case err := <-receiveStream:
			if err != nil {
				return err
			}
		}
	}

	status := <-statusChan
	client.ContainerRemove(ctx, res.ID, types.ContainerRemoveOptions{})
	return s.Exit(int(status))
}
