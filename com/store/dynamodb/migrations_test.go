package dynamodb

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestMigrateApply(t *testing.T) {
	testMigrations := Migrations{
		1: {
			version:     1,
			description: "rename config to env",
			apply: func(in map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
				return in, fmt.Errorf("err")
			},
		},
	}
	cmd := map[string]*dynamodb.AttributeValue{
		"Config": &dynamodb.AttributeValue{
			M: map[string]*dynamodb.AttributeValue{
				"name": {S: aws.String("var")},
			},
		},
	}

	_, err := testMigrations.Apply(latestVesion, cmd)
	assert.Error(t, err,
		"expected apply to fail with error about missing version")

	cmd[schemaAttr] = &dynamodb.AttributeValue{
		N: aws.String("0"),
	}
	res, err := testMigrations.Apply(latestVesion, cmd)
	assert.Error(t, err,
		"expected apply to fail with err = \"err\"")
	assert.NotNil(t, res)
}

func TestAllMigrations(t *testing.T) {
	cmd := map[string]*dynamodb.AttributeValue{
		"Config": &dynamodb.AttributeValue{
			M: map[string]*dynamodb.AttributeValue{
				"name": {S: aws.String("var")},
			},
		},
		schemaAttr: &dynamodb.AttributeValue{
			N: aws.String("0"),
		},
	}

	res, err := migrations.Apply(latestVesion, cmd)
	assert.NoError(t, err,
		"failed to apply migrations")
	assert.NotNil(t, res)

	assert.NotEqual(t, cmd, res,
		"migrations must return copy of item")

	assert.Contains(t, res, "Environment")
	assert.NotNil(t, res["Environment"])

	assert.Equal(t, cmd["Config"].M["name"], res["Environment"].M["name"])

	assert.NotContains(t, res, "Config",
		"config attribute should have been removed in migration")

	assert.Equal(t, latestVesion, itemVersion(res))
}
