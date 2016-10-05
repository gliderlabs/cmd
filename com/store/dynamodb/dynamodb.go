package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	dynamoattr "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gliderlabs/gosper/pkg/com"
	"github.com/gliderlabs/gosper/pkg/log"
	"github.com/progrium/cmd/com/core"
)

func init() {
	com.Register("store.dynamodb", &Component{},
		com.Option("table", "", "dynamodb table name for command storage"),
		com.Option("access_key", "", "aws access key for dynamodb store"),
		com.Option("secret_key", "", "aws secret key for dynamodb store"),
		com.Option("region", "us-east-1", "aws region for dynamodb store"),
	)
}

type Component struct{}

func (c *Component) client() *dynamodb.DynamoDB {
	return dynamodb.New(session.New(
		aws.NewConfig().
			WithRegion(
				com.GetString("region")).
			WithCredentials(credentials.NewStaticCredentials(
				com.GetString("access_key"),
				com.GetString("secret_key"), ""))))
}

func (c *Component) List(user string) []*core.Command {
	res, err := c.client().Scan(&dynamodb.ScanInput{
		ScanFilter: map[string]*dynamodb.Condition{
			"User": &dynamodb.Condition{
				AttributeValueList: []*dynamodb.AttributeValue{
					&dynamodb.AttributeValue{S: aws.String(user)},
				},
				ComparisonOperator: aws.String("EQ"),
			},
		},
		TableName: aws.String(com.GetString("table")),
	})
	if err != nil {
		// FIXME: Should actually do something with this error.
		log.Info(err)
		return nil
	}

	cmds := make([]*core.Command, 0, len(res.Items))
	for _, item := range res.Items {
		var cmd core.Command
		err := dynamoattr.UnmarshalMap(item, &cmd)
		if err != nil {
			log.Debug(err)
			continue
		}
		cmds = append(cmds, &cmd)
	}
	return cmds
}

func (c *Component) Get(user, name string) *core.Command {
	res, err := c.client().GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"User": &dynamodb.AttributeValue{
				S: aws.String(user),
			},
			"Name": &dynamodb.AttributeValue{
				S: aws.String(name),
			},
		},
		TableName: aws.String(com.GetString("table")),
	})
	if err != nil {
		log.Info(err)
		return nil
	}

	if res.Item == nil {
		return nil
	}

	var cmd core.Command
	if err = dynamoattr.UnmarshalMap(res.Item, &cmd); err != nil {
		log.Debug(err)
	}
	return &cmd
}

func (c *Component) Put(user, name string, cmd *core.Command) error {
	item, err := dynamoattr.MarshalMap(cmd)
	if err != nil {
		return err
	}

	_, err = c.client().PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(com.GetString("table")),
	})
	return err
}

func (c *Component) Delete(user, name string) error {
	_, err := c.client().DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"User": &dynamodb.AttributeValue{
				S: aws.String(user),
			},
			"Name": &dynamodb.AttributeValue{
				S: aws.String(name),
			},
		},
		TableName: aws.String(com.GetString("table")),
	})
	return err
}
