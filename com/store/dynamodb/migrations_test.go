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

	actual := assertMigrationApply(t, 1, cmd)
	assert.Contains(t, actual, "Environment")
	assert.NotNil(t, actual["Environment"])
	assert.Equal(t,
		cmd["Config"].M["name"],
		actual["Environment"].M["name"],
		"Environment should match Config")

	assert.NotContains(t, actual, "Config",
		"config attribute should have been removed in migration")
}

func TestSchemaV2(t *testing.T) {
	cmd := map[string]*dynamodb.AttributeValue{
		"ACL": &dynamodb.AttributeValue{
			L: []*dynamodb.AttributeValue{
				{S: aws.String("*")},
			},
		},
		schemaAttr: &dynamodb.AttributeValue{
			N: aws.String("1"),
		},
	}

	actual := assertMigrationApply(t, 2, cmd)
	assert.Equal(t,
		&dynamodb.AttributeValue{NULL: aws.Bool(true)},
		actual["ACL"],
		"NULL AttributeValue should replace public ACL")

	cmd = map[string]*dynamodb.AttributeValue{
		"ACL": &dynamodb.AttributeValue{
			L: []*dynamodb.AttributeValue{
				{S: aws.String("user")},
			},
		},
		schemaAttr: &dynamodb.AttributeValue{
			N: aws.String("1"),
		},
	}

	actual = assertMigrationApply(t, 2, cmd)
	assert.Equal(t, cmd["ACL"], actual["ACL"],
		"No modification should be made unless ACL contains *")

	cmd = map[string]*dynamodb.AttributeValue{
		"ACL": &dynamodb.AttributeValue{
			L: []*dynamodb.AttributeValue{
				{S: aws.String("user")},
				{S: aws.String("*")},
			},
		},
		schemaAttr: &dynamodb.AttributeValue{
			N: aws.String("1"),
		},
	}

	actual = assertMigrationApply(t, 2, cmd)
	assert.Equal(t,
		&dynamodb.AttributeValue{NULL: aws.Bool(true)},
		actual["ACL"],
		"NULL AttributeValue should replace list containing public ACL")
}

func assertMigrationApply(t *testing.T, target int, item map[string]*dynamodb.AttributeValue) map[string]*dynamodb.AttributeValue {
	res, err := migrations.Apply(target, item)
	assert.NoError(t, err,
		"failed to apply migrations")
	assert.NotNil(t, res)
	assert.NotEqual(t, item, res,
		"migrations must return copy of item")
	assert.Equal(t, target, itemVersion(res))

	return res
}
