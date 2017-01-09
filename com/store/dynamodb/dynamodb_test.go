package dynamodb

import (
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gliderlabs/gosper/pkg/com"
	"github.com/gliderlabs/gosper/pkg/com/viper"
	"github.com/progrium/cmd/com/core"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	os.Setenv("DYNAMODB_TABLE", "cmd-test-table")
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:8000")
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

}

// ensureTableExists creates a DynamoDB table with a given
// DynamoDB client. If the table already exists, it is not
// being reconfigured.
func ensureTableExists(client *dynamodb.DynamoDB, table string, readCapacity, writeCapacity int) error {
	_, err := client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})
	if awserr, ok := err.(awserr.Error); ok {
		if awserr.Code() == "ResourceNotFoundException" {
			_, err = client.CreateTable(&dynamodb.CreateTableInput{
				TableName: aws.String(table),
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(int64(readCapacity)),
					WriteCapacityUnits: aws.Int64(int64(writeCapacity)),
				},
				KeySchema: []*dynamodb.KeySchemaElement{{
					AttributeName: aws.String("User"),
					KeyType:       aws.String("HASH"),
				}, {
					AttributeName: aws.String("Name"),
					KeyType:       aws.String("RANGE"),
				}},
				AttributeDefinitions: []*dynamodb.AttributeDefinition{{
					AttributeName: aws.String("User"),
					AttributeType: aws.String("S"),
				}, {
					AttributeName: aws.String("Name"),
					AttributeType: aws.String("S"),
				}},
			})
			if err != nil {
				return err
			}
			err = client.WaitUntilTableExists(&dynamodb.DescribeTableInput{
				TableName: aws.String(table),
			})
			if err != nil {
				return err
			}
		}
	}

	return err
}
