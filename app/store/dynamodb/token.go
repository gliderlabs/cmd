package dynamodb

import (
	"github.com/gliderlabs/cmd/app/core"
)

// ListTokens for a user.
func (c *Component) ListTokens(user string) ([]*core.Token, error) {
	var tokens []*core.Token
	err := c.tokenTable().Scan().Filter("'User' = ?", user).All(&tokens)
	return tokens, err
}

// GetToken by id.
func (c *Component) GetToken(id string) (*core.Token, error) {
	var token *core.Token
	err := c.tokenTable().Get("Key", id).One(&token)
	return token, err
}

// PutToken ...
func (c *Component) PutToken(token *core.Token) error {
	return c.tokenTable().Put(token).Run()
}

// DeleteToken by id
func (c *Component) DeleteToken(id string) error {
	return c.tokenTable().Delete("Key", id).Run()
}
