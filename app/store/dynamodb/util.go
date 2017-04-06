package dynamodb

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gliderlabs/comlab/pkg/log"
	maintenance "github.com/gliderlabs/cmd/lib/maint"
)

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
			if aws.StringValue(client.Config.Region) == "local" {
				log.Info("skipping unsupported dynamodb-local operation", log.Fields{"operation": "setTableVersion"})
				return nil
			}
			return setTableVersion(client, table, latestVesion)
		}
	}

	return err
}

// ensureTokenTableExists creates a DynamoDB table with a given
// DynamoDB client. If the table already exists, it is not
// being reconfigured.
func ensureTokenTableExists(client *dynamodb.DynamoDB, table string, readCapacity, writeCapacity int) error {
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
					AttributeName: aws.String("Key"),
					KeyType:       aws.String("HASH"),
				}},
				AttributeDefinitions: []*dynamodb.AttributeDefinition{{
					AttributeName: aws.String("Key"),
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

func setTableVersion(client *dynamodb.DynamoDB, name string, version int) error {
	arn := tableArn(client, name)
	_, err := client.TagResource(&dynamodb.TagResourceInput{
		ResourceArn: aws.String(arn),
		Tags: []*dynamodb.Tag{{
			Key:   aws.String(tableSchemaKey),
			Value: aws.String(strconv.Itoa(version)),
		}},
	})
	return err
}

func getTableVersion(client *dynamodb.DynamoDB, name string) int {
	arn := tableArn(client, name)
	res, err := client.ListTagsOfResource(&dynamodb.ListTagsOfResourceInput{
		ResourceArn: aws.String(arn),
	})
	if err != nil {
		return -1
	}

	for _, tag := range res.Tags {
		if *tag.Key == tableSchemaKey {
			version, _ := strconv.Atoi(*tag.Value)
			return version
		}
	}
	return -1
}

func tableArn(client *dynamodb.DynamoDB, name string) string {
	res, err := client.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(name),
	})
	if err != nil {
		log.Debug(err)
		return ""
	}
	return aws.StringValue(res.Table.TableArn)
}

func containsArg(args []string, s string) bool {
	for i := 0; i < len(args); i++ {
		if args[i] == s {
			return true
		}
	}
	return false
}

func ensureTableSchema(client *dynamodb.DynamoDB, table string) error {
	if containsArg(os.Args, "-migrate") {
		if !maintenance.Active() {
			return errors.New("maintenance must be active when passing -migrate flag")
		}

		if err := migrateTable(client, table, latestVesion); err != nil {
			return err
		}
		fmt.Println("done")
		os.Exit(0)
	}

	if aws.StringValue(client.Config.Region) == "local" {
		log.Info("applying migrations to dynamodb-local")
		if err := migrateTable(client, table, latestVesion); err != nil {
			return err
		}
	} else {
		current := getTableVersion(client, table)
		if !maintenance.Active() && migrations.HardRequired(current, latestVesion) {
			return errors.New("hard migration required")
		}
	}
	return nil
}

func stringSet(item ...string) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{SS: aws.StringSlice(item)}
}
