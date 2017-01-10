package dynamodb

import (
	"bytes"
	"encoding/gob"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

const (
	schemaAttr   = "_schema"
	latestVesion = 2
)

type migration struct {
	version     int
	description string

	apply func(map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error)
}

var migrations = Migrations{
	1: {
		version:     1,
		description: "rename config to env",
		apply: func(in map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
			out := map[string]*dynamodb.AttributeValue{}
			err := copyAttrVal(out, in)
			if err != nil {
				return out, err
			}
			if _, ok := in["Environment"]; !ok {
				out["Environment"] = in["Config"]
				delete(out, "Config")
			}
			return out, nil
		},
	},
	2: {
		version:     2,
		description: "drop public cmd support",
		apply: func(in map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
			out := map[string]*dynamodb.AttributeValue{}
			err := copyAttrVal(out, in)
			if err != nil {
				return out, err
			}
			if acl := in["ACL"]; acl != nil {
				for _, v := range acl.L {
					if aws.StringValue(v.S) == "*" {
						out["ACL"] = &dynamodb.AttributeValue{
							NULL: aws.Bool(true),
						}
						return out, nil
					}
				}
			}
			return out, nil
		},
	},
}

// Migrations contains a map of migrations keyed by the version number
type Migrations map[int]migration

// Apply all migrations on copy of item until item version is equal to target.
func (ms Migrations) Apply(target int, item map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	initialVer := itemVersion(item)
	cur := item
	for nextVer := initialVer + 1; nextVer <= target; nextVer++ {
		m, ok := ms[nextVer]
		if !ok {
			continue
		}
		var next map[string]*dynamodb.AttributeValue
		next, err := m.apply(cur)
		if err != nil {
			return cur, errors.Wrapf(err,
				"unable to apply migration version: %d - %s", nextVer, m.description)
		}
		next[schemaAttr] = &dynamodb.AttributeValue{
			N: aws.String(strconv.Itoa(nextVer)),
		}
		cur = next
	}
	return cur, nil
}

func itemVersion(item map[string]*dynamodb.AttributeValue) int {
	if _, ok := item[schemaAttr]; !ok {
		return 0
	}
	v, _ := dynamodbattribute.Number(aws.StringValue(item[schemaAttr].N)).Int64()
	return int(v)
}

func init() {
	gob.Register(&map[string]*dynamodb.AttributeValue{})
}

// copyAttrVal performs a deep copy of the given attr map map src.
func copyAttrVal(dst, src map[string]*dynamodb.AttributeValue) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(src)
	if err != nil {
		return err
	}
	return dec.Decode(&dst)
}
