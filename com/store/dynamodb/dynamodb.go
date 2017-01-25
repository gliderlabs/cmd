package dynamodb

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	dynamoattr "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/pkg/errors"
	"github.com/progrium/cmd/com/core"
)

func init() {
	com.Register("store.dynamodb", &Component{},
		com.Option("table", "", "dynamodb table name for command storage"),
		com.Option("token_table", "", "dynamodb table name for token storage"),
		com.Option("access_key", "", "aws access key for dynamodb store"),
		com.Option("secret_key", "", "aws secret key for dynamodb store"),
		com.Option("region", "us-east-1", "aws region for dynamodb store"),
		com.Option("endpoint", "", "alternate dynamodb endpoint. eg: http://localhost:8000"),
	)
}

type Component struct{}

func (c *Component) AppPreStart() error {
	err := ensureTableExists(c.client(), com.GetString("table"), 5, 5)
	if err != nil {
		return err
	}
	return ensureTokenTableExists(c.client(), com.GetString("token_table"), 5, 5)
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
		log.Info(errors.Wrapf(err, "unable to list commands for user: %s", user))
		return nil
	}

	cmds := make([]*core.Command, 0, len(res.Items))
	for _, item := range res.Items {
		migrated, err := migrations.Apply(latestVesion, item)
		if err != nil {
			log.Info(errors.Wrapf(err,
				"failed migrating commands for user: %s to version: %d",
				user, latestVesion),
			)
			continue
		}
		var cmd core.Command
		if err := dynamoattr.UnmarshalMap(migrated, &cmd); err != nil {
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
		log.Info(errors.Wrapf(err, "unable to get cmd: %s for user: %s", name, user))
		return nil
	}

	if res.Item == nil {
		return nil
	}
	migrated, err := migrations.Apply(latestVesion, res.Item)
	if err != nil {
		log.Info(errors.Wrapf(err,
			"failed migrating cmd: %s from user: %s to version: %s",
			name, user, latestVesion),
		)
		return nil
	}
	var cmd core.Command
	if err = dynamoattr.UnmarshalMap(migrated, &cmd); err != nil {
		log.Debug(err)
	}
	return &cmd
}

func (c *Component) Put(user, name string, cmd *core.Command) error {
	item, err := dynamoattr.MarshalMap(cmd)
	if err != nil {
		return err
	}

	if _, ok := item[schemaAttr]; !ok {
		item[schemaAttr] = &dynamodb.AttributeValue{
			N: aws.String(strconv.Itoa(latestVesion))}
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

func (c *Component) ListTokens(user string) ([]*core.Token, error) {
	res, err := c.client().Scan(&dynamodb.ScanInput{
		ScanFilter: map[string]*dynamodb.Condition{
			"User": &dynamodb.Condition{
				AttributeValueList: []*dynamodb.AttributeValue{
					&dynamodb.AttributeValue{S: aws.String(user)},
				},
				ComparisonOperator: aws.String("EQ"),
			},
		},
		TableName: aws.String(com.GetString("token_table")),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to list tokens for user:%s", user)
	}

	tokens := make([]*core.Token, 0, len(res.Items))
	for _, item := range res.Items {
		var token core.Token
		if err := dynamoattr.UnmarshalMap(item, &token); err != nil {
			log.Debug(err)
			continue
		}
		tokens = append(tokens, &token)
	}
	return tokens, nil
}

func (c *Component) GetToken(id string) (*core.Token, error) {
	res, err := c.client().GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Key": &dynamodb.AttributeValue{
				S: aws.String(id),
			},
		},
		TableName: aws.String(com.GetString("token_table")),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get token: %s", id)
	}
	if res.Item == nil {
		return nil, nil
	}
	var token core.Token
	return &token, dynamoattr.UnmarshalMap(res.Item, &token)
}

func (c *Component) PutToken(token *core.Token) error {
	item, err := dynamoattr.MarshalMap(token)
	if err != nil {
		return err
	}

	_, err = c.client().PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(com.GetString("token_table")),
	})
	return err
}

func (c *Component) DeleteToken(id string) error {
	_, err := c.client().DeleteItem(&dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Key": &dynamodb.AttributeValue{
				S: aws.String(id),
			},
		},
		TableName: aws.String(com.GetString("token_table")),
	})
	return err
}
