package dynamodb

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	dynamoattr "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
	"github.com/gliderlabs/cmd/app/core"
)

// List all commands for a given user.
func (c *Component) List(user string) []*core.Command {
	var items []map[string]*dynamodb.AttributeValue
	err := c.cmdTable().Scan().Filter("'User' = ?", user).All(&items)
	if err != nil {
		log.Info(errors.Wrapf(err, "unable to list commands for user: %s", user))
		return nil
	}
	cmds := make([]*core.Command, 0, len(items))
	for _, item := range items {
		migrated, err := migrations.Apply(latestVesion, false, item)
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

// Get command by name for given user
func (c *Component) Get(user, name string) *core.Command {
	item, err := c.cmdTable().Get("User", user).Range("Name", dynamo.Equal, name).OneItem()
	if err != nil {
		return nil
	}
	migrated, err := migrations.Apply(latestVesion, false, item)
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

// Put command with name for given user.
func (c *Component) Put(user, name string, cmd *core.Command) error {
	item, err := dynamoattr.MarshalMap(cmd)
	if err != nil {
		return err
	}

	if _, ok := item[schemaAttr]; !ok {
		item[schemaAttr] = &dynamodb.AttributeValue{
			N: aws.String(strconv.Itoa(latestVesion))}
	}
	return c.cmdTable().Put(item).Run()
}

// Delete command by name  for given user
func (c *Component) Delete(user, name string) error {
	return c.cmdTable().
		Delete("User", user).Range("Name", name).Run()
}

func (c *Component) updateCmd(owner, name string) *dynamo.Update {
	return c.cmdTable().
		Update("User", owner).Range("Name", name)
}

// GrantAccess to a command, for each subject
func (c *Component) GrantAccess(owner, name string, subject ...string) error {
	return c.updateCmd(owner, name).
		Add("ACL", stringSet(subject...)).Run()
}

// RevokeAccess to a command, for each subject
func (c *Component) RevokeAccess(owner, name string, subject ...string) error {
	return c.updateCmd(owner, name).
		Delete("ACL", stringSet(subject...)).Run()
}

// GrantAdmin to a command, for each subject
func (c *Component) GrantAdmin(owner, name string, subject ...string) error {
	return c.updateCmd(owner, name).
		Add("Admins", stringSet(subject...)).Run()
}

// RevokeAdmin to a command, for each subject
func (c *Component) RevokeAdmin(owner, name string, subject ...string) error {
	return c.updateCmd(owner, name).
		Delete("Admins", stringSet(subject...)).Run()
}
