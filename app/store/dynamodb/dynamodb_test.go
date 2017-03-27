package dynamodb

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/com/viper"
	"github.com/stretchr/testify/assert"

	"github.com/progrium/cmd/app/core"
	"github.com/progrium/cmd/app/store"
)

func TestStore(t *testing.T) {
	assert.Implements(t, new(store.Backend), new(Component))

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
		err := c.Put("user", "cmd", &core.Command{
			User: "user",
			Name: "cmd",
		})
		assert.NoError(t, err)
	})

	t.Run("PutNil", func(t *testing.T) {
		err := c.Put("user", "cmd", nil)
		assert.Error(t, err)
	})

	t.Run("Get", func(t *testing.T) {
		assert.NotNil(t, c.Get("user", "cmd"))
	})

	t.Run("List", func(t *testing.T) {
		assert.NotNil(t, c.List("user"))
	})
	t.Run("Delete", func(t *testing.T) {
		assert.NoError(t, c.Delete("user", "cmd"))
		assert.Nil(t, c.Get("user", "cmd"))
	})

	ensureTokenTableExists(c.client(), "cmd-test-tokens-table", 5, 5)
	t.Run("PutToken", func(t *testing.T) {
		assert.NoError(t, c.PutToken(&core.Token{Key: "key", User: "user"}))
	})

	t.Run("GetToken", func(t *testing.T) {
		token, err := c.GetToken("key")
		assert.NoError(t, err)
		assert.NotNil(t, token)
	})

	t.Run("ListTokens", func(t *testing.T) {
		tokens, err := c.ListTokens("user")
		assert.NoError(t, err)
		if assert.NotNil(t, tokens) {
			assert.Len(t, tokens, 1)
		}
	})

	t.Run("DeleteToken", func(t *testing.T) {
		assert.NoError(t, c.DeleteToken("key"))
		token, err := c.GetToken("key")
		assert.NoError(t, err)
		assert.Nil(t, token)
	})
}
