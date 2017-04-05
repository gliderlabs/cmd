package store

import (
	"github.com/progrium/cmd/app/core"

	"github.com/gliderlabs/comlab/pkg/com"
)

func init() {
	com.Register("store", struct{}{},
		com.Option("backend", "store.filesystem", "Store backend"))
}

func Selected() Backend {
	backend := com.Select(com.GetString("backend"), new(Backend))
	if backend == nil {
		panic("Unable to find selected backend: " + com.GetString("backend"))
	}
	return backend.(Backend)
}

type Backend interface {
	CmdBackend
	TokenBackend
}

type CmdBackend interface {
	List(user string) []*core.Command
	Get(user, name string) *core.Command
	Put(user, name string, cmd *core.Command) error
	Delete(user, name string) error
	GrantAccess(owner, name string, subject ...string) error
	RevokeAccess(owner, name string, subject ...string) error
	GrantAdmin(owner, name string, subject ...string) error
	RevokeAdmin(owner, name string, subject ...string) error
}

type TokenBackend interface {
	ListTokens(user string) ([]*core.Token, error)
	GetToken(key string) (*core.Token, error)
	PutToken(token *core.Token) error
	DeleteToken(key string) error
}
