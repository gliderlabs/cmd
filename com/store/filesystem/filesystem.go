package filesystem

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"

	"github.com/progrium/cmd/com/core"
)

func init() {
	com.Register("store.filesystem", &Component{},
		com.Option("dir", "local", "directory for file store"))
}

type Component struct{}

func (c *Component) load(filepath string) *core.Command {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Debug(err)
		return nil
	}
	var cmd core.Command
	if err := json.Unmarshal(b, &cmd); err == nil {
		return &cmd
	} else {
		log.Debug(err)
		return nil
	}
}

func (c *Component) List(user string) []*core.Command {
	matches, _ := filepath.Glob(com.GetString("dir") + "/" + user + "-*")
	var cmds []*core.Command
	for _, match := range matches {
		if cmd := c.load(match); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}

func (c *Component) path(user, name string) string {
	return com.GetString("dir") + "/" + user + "-" + name + ".json"
}

func (c *Component) Get(user, name string) *core.Command {
	return c.load(c.path(user, name))
}

func (c *Component) Put(user, name string, cmd *core.Command) error {
	b, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.path(user, name), b, 0644)
}

func (c *Component) Delete(user, name string) error {
	return os.Remove(c.path(user, name))
}
