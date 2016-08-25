package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gliderlabs/pkg/com"
	"github.com/gliderlabs/pkg/log"
)

type CommandStore interface {
	List(user string) []*Command
	Get(user, name string) *Command
	Put(user, name string, cmd *Command) error
	Delete(user, name string) error
}

func GetStore() CommandStore {
	return &fileStore{}
}

type fileStore struct{}

func (s *fileStore) load(filepath string) *Command {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Debug(err)
		return nil
	}
	var cmd Command
	if err := json.Unmarshal(b, &cmd); err == nil {
		return &cmd
	} else {
		log.Debug(err)
		return nil
	}
}

func (s *fileStore) List(user string) []*Command {
	matches, _ := filepath.Glob(com.GetString("config_dir") + "/" + user + "-*")
	var cmds []*Command
	for _, match := range matches {
		if cmd := s.load(match); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}

func (s *fileStore) path(user, name string) string {
	return com.GetString("config_dir") + "/" + user + "-" + name + ".json"
}

func (s *fileStore) Get(user, name string) *Command {
	return s.load(s.path(user, name))
}

func (s *fileStore) Put(user, name string, cmd *Command) error {
	b, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(s.path(user, name), b, 0644)
}

func (s *fileStore) Delete(user, name string) error {
	return os.Remove(s.path(user, name))
}
