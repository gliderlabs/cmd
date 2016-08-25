package cmd

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/gliderlabs/pkg/com"
	"github.com/gliderlabs/pkg/ssh"
)

type Command struct {
	Name        string
	User        string
	Config      map[string]string
	ACL         []string
	Admins      []string
	Description string
	Source      string
}

func (c *Command) IsPublic() bool {
	return len(c.ACL) == 1 && c.ACL[0] == "*"
}

func (c *Command) MakePublic() {
	c.ACL = []string{"*"}
}

func (c *Command) MakePrivate() {
	c.ACL = []string{}
}

func (c *Command) HasAccess(user string) bool {
	if c.IsPublic() {
		return true
	}
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

func (c *Command) AddAccess(user string) {
	if c.HasAccess(user) {
		return
	}
	c.ACL = append(c.ACL, user)
}

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

func (c *Command) Pull() error {
	pull := exec.Command(com.GetString("docker_bin"), "pull", c.Source)
	if err := pull.Run(); err != nil {
		return err
	}
	tag := exec.Command(com.GetString("docker_bin"), "tag", c.Source, c.image())
	return tag.Run()
}

func (c *Command) UpdateDescription() error {
	descCmd := exec.Command(com.GetString("docker_bin"),
		[]string{"inspect", "-f", `{{ index .ContainerConfig.Labels "io.cmd.description" }}`, c.image()}...)
	if b, err := descCmd.Output(); err != nil {
		return err
	} else {
		c.Description = strings.Trim(string(b), "\n")
		if err := Store.Put(c.User, c.Name, c); err != nil {
			return err
		}
	}
	return nil
}

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
	runArgs := []string{"run", "--rm", "-i"}
	for k, v := range c.Config {
		runArgs = append(runArgs, "-e")
		runArgs = append(runArgs, fmt.Sprintf("%s=%s", k, v))
	}
	runArgs = append(runArgs, c.image())
	docker := exec.Command(com.GetString("docker_bin"), append(runArgs, args...)...)
	docker.Stdout = s
	docker.Stderr = s.Stderr()
	stdinPipe, err := docker.StdinPipe()
	if err != nil {
		fmt.Fprintln(s.Stderr(), err.Error())
		s.Exit(255)
		return
	}
	go func() {
		defer stdinPipe.Close()
		io.Copy(stdinPipe, s)
	}()
	if err := docker.Run(); err != nil {
		s.Exit(1)
	}
}
