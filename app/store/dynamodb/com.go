package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

func init() {
	com.Register("store.dynamodb", &Component{},
		com.Option("table", "", "dynamodb table name for command storage"),
		com.Option("token_table", "", "dynamodb table name for token storage"),
		com.Option("access_key", "", "aws access key for dynamodb store"),
		com.Option("secret_key", "", "aws secret key for dynamodb store"),
		com.Option("endpoint", "", "alternate dynamodb endpoint. eg: http://localhost:8000"),
		com.Option("region", "us-east-1", "aws region for dynamodb store, NOTE: use 'local' for dynamodb-local"),
	)
}

// Component implements a store backend
type Component struct{}

// AppPreStart attempts to create any missing dynamodb tables and ensures
// cmd table has the latest schema, optionally applying migrations when
// maintenance is active and `-migrate` flag passed.
//
// Note that an error will be returned if a hard migration is required and
// maintenance is NOT active.
func (c *Component) AppPreStart() error {
	var (
		cmdTable   = com.GetString("table")
		tokenTable = com.GetString("token_table")
	)

	if err := ensureTableExists(c.client(), cmdTable, 5, 5); err != nil {
		return errors.Wrapf(err, "dynamodb table %q setup failed", cmdTable)
	}

	if err := ensureTokenTableExists(c.client(), tokenTable, 5, 5); err != nil {
		return errors.Wrapf(err, "dynamodb table %q setup failed", tokenTable)
	}

	return ensureTableSchema(c.client(), cmdTable)
}

func (c *Component) cmdTable() dynamo.Table {
	db := dynamo.New(session.New(), &c.client().Config)
	return db.Table(com.GetString("table"))
}

func (c *Component) tokenTable() dynamo.Table {
	db := dynamo.New(session.New(), &c.client().Config)
	return db.Table(com.GetString("token_table"))
}

func (c *Component) client() *dynamodb.DynamoDB {
	var (
		region    = com.GetString("region")
		accessKey = com.GetString("access_key")
		secretKey = com.GetString("secret_key")
		endpoint  = com.GetString("endpoint")
		retries   = com.GetInt("max_retries")
	)
	config := aws.NewConfig().
		WithRegion(region).
		WithCredentials(
			credentials.NewStaticCredentials(
				accessKey, secretKey, ""),
		)

	if endpoint != "" {
		config = config.WithEndpoint(endpoint)
	}

	if retries != 0 {
		config = config.WithMaxRetries(retries)
	}
	return dynamodb.New(session.New(config))
}
