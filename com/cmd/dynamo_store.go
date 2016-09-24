package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
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

func GetDynamodbStore(sess *session.Session) CommandStore {
	return &dynamodbStore{
		table:  com.GetString("TABLE_NAME"),
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

	var cmds []*Command
	for _, item := range res.Items {
		var cmd Command
		err := dynamoattr.UnmarshalMap(item, &cmd)
		if err != nil {
			log.Debug(err)
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
