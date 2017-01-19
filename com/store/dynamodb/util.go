package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
