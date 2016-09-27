package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	dynamoattr "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gliderlabs/pkg/com"
	"github.com/gliderlabs/pkg/log"
)

type dynamodbStore struct {
	table string

	client *dynamodb.DynamoDB
}

func GetDynamodbStore() CommandStore {
	sess := session.New(
		aws.NewConfig().
			WithRegion(
				com.GetString("aws_region")).
			WithCredentials(credentials.NewStaticCredentials(
				com.GetString("aws_access_key"),
				com.GetString("aws_secret_key"), "")))

	return &dynamodbStore{
		table:  com.GetString("table_name"),
		client: dynamodb.New(sess),
	}
}

func (s *dynamodbStore) List(user string) []*Command {
	res, err := s.client.Scan(&dynamodb.ScanInput{
		ScanFilter: map[string]*dynamodb.Condition{
			"User": &dynamodb.Condition{
				AttributeValueList: []*dynamodb.AttributeValue{
					&dynamodb.AttributeValue{S: aws.String(user)},
				},
				ComparisonOperator: aws.String("EQ"),
			},
		},
		TableName: aws.String(s.table),
	})
	if err != nil {
		// FIXME: Should actually do something with this error.
		log.Debug(err)
		return nil
	}

	cmds := make([]*Command, 0, len(res.Items))
	for _, item := range res.Items {
		var cmd Command
		err := dynamoattr.UnmarshalMap(item, &cmd)
		if err != nil {
			log.Debug(err)
			continue
		}
		cmds = append(cmds, &cmd)
	}
	return cmds
}

func (s *dynamodbStore) Get(user, name string) *Command {
	res, err := s.client.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"User": &dynamodb.AttributeValue{
				S: aws.String(user),
			},
			"Name": &dynamodb.AttributeValue{
				S: aws.String(name),
			},
		},
		TableName: aws.String(s.table),
	})
	if err != nil {
		log.Debug(err)
		return nil
	}

	var cmd Command
	if err = dynamoattr.UnmarshalMap(res.Item, &cmd); err != nil {
		log.Debug(err)
	}
	return &cmd
}

func (s *dynamodbStore) Put(user, name string, cmd *Command) error {
	item, err := dynamoattr.MarshalMap(cmd)
	if err != nil {
		return err
	}

	_, err = s.client.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(s.table),
	})
	return err
}

func (s *dynamodbStore) Delete(user, name string) error {
	_, err := s.client.DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"User": &dynamodb.AttributeValue{
				S: aws.String(user),
			},
			"Name": &dynamodb.AttributeValue{
				S: aws.String(name),
			},
		},
		TableName: aws.String(s.table),
	})
	return err
}
