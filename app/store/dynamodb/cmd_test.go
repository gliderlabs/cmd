package dynamodb

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/com/viper"
	"github.com/stretchr/testify/assert"

	"github.com/gliderlabs/cmd/app/core"
	"github.com/gliderlabs/cmd/app/store"
)

func TestCmdBackend(t *testing.T) {
	assert.Implements(t, new(store.CmdBackend), new(Component))

	os.Setenv("DYNAMODB_TABLE", "cmd-test-table")
	os.Setenv("DYNAMODB_TOKEN_TABLE", "cmd-test-tokens-table")
	os.Setenv("DYNAMODB_REGION", "local")
	os.Setenv("DYNAMODB_ACCESS_KEY", "test")
	os.Setenv("DYNAMODB_SECRET_KEY", "test")
	os.Setenv("DYNAMODB_MAX_RETRIES", "1")
	cfg := viper.NewConfig()
	cfg.AutomaticEnv()
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	com.SetConfig(cfg)
	c := &Component{}

	err := ensureTableExists(c.client(), "cmd-test-table", 5, 5)
	if awserr, ok := err.(awserr.Error); ok {
		if awserr.Code() == "RequestError" && awserr.Message() == "send request failed" {
			t.Skip("unable to connect to local instance of dynamodb", awserr)
		}
	}

	t.Run("Put", func(t *testing.T) {
		assert.NoError(t, c.Put("user", "cmd", &core.Command{
			User: "user",
			Name: "cmd",
		}))
	})

	t.Run("PutNil", func(t *testing.T) {
		assert.Error(t, c.Put("user", "cmd", nil))
	})

	t.Run("Get", func(t *testing.T) {
		assert.NotNil(t, c.Get("user", "cmd"))
	})

	t.Run("List", func(t *testing.T) {
		assert.NoError(t, c.Put("user2", "cmd", &core.Command{
			User: "user2",
			Name: "cmd",
		}))

		cmds := c.List("user")
		assert.NotNil(t, cmds)
		for _, cmd := range cmds {
			assert.Equal(t, "user", cmd.User,
				"Result should only contain commands owned by 'user'")
		}
	})

	t.Run("GrantAccess", func(t *testing.T) {
		err := c.GrantAccess("user", "cmd", "foo")
		assert.NoError(t, err)
		cmd := c.Get("user", "cmd")
		if assert.NotNil(t, cmd) {
			assert.EqualValues(t, []string{"foo"}, cmd.ACL)
		}
		err = c.GrantAccess("user", "cmd", "bar")
		assert.NoError(t, err)
		cmd = c.Get("user", "cmd")
		if assert.NotNil(t, cmd) {
			assert.EqualValues(t, []string{"bar", "foo"}, cmd.ACL)
		}
	})

	t.Run("RevokeAccess", func(t *testing.T) {
		err := c.RevokeAccess("user", "cmd", "foo")
		assert.NoError(t, err)
		cmd := c.Get("user", "cmd")
		if assert.NotNil(t, cmd) {
			assert.EqualValues(t, []string{"bar"}, cmd.ACL)
		}
	})

	t.Run("GrantAdmin", func(t *testing.T) {
		err := c.GrantAdmin("user", "cmd", "foo")
		assert.NoError(t, err)
		cmd := c.Get("user", "cmd")
		if assert.NotNil(t, cmd) {
			assert.EqualValues(t, []string{"foo"}, cmd.Admins)
		}
		err = c.GrantAdmin("user", "cmd", "bar")
		assert.NoError(t, err)
		cmd = c.Get("user", "cmd")
		if assert.NotNil(t, cmd) {
			assert.EqualValues(t, []string{"bar", "foo"}, cmd.Admins)
		}
	})

	t.Run("RevokeAdmin", func(t *testing.T) {
		err := c.RevokeAdmin("user", "cmd", "foo")
		assert.NoError(t, err)
		cmd := c.Get("user", "cmd")
		if assert.NotNil(t, cmd) {
			assert.EqualValues(t, []string{"bar"}, cmd.Admins)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		assert.NoError(t, c.Delete("user", "cmd"))
		assert.Nil(t, c.Get("user", "cmd"))
	})
}
