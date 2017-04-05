package dynamodb

import (
	"os"
	"strings"
	"testing"

	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/com/viper"
	"github.com/guregu/dynamo"
	"github.com/stretchr/testify/assert"

	"github.com/progrium/cmd/app/core"
	"github.com/progrium/cmd/app/store"
)

func TestTokenBackend(t *testing.T) {
	assert.Implements(t, new(store.TokenBackend), new(Component))

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
		assert.EqualError(t, err, dynamo.ErrNotFound.Error())
		assert.Nil(t, token)
	})
}
