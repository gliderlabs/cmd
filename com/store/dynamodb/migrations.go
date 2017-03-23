package dynamodb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cenkalti/backoff"
	"github.com/gliderlabs/comlab/pkg/log"
	"github.com/pkg/errors"
)

const (
	schemaAttr     = "_schema"
	tableSchemaKey = "cmd:_schema"
	latestVesion   = 3

	// maxBatchSize dynamodb allows in a single request
	maxBatchSize = 25
)

type migration struct {
	version     int
	description string
	hard        bool

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
	3: {
		version:     3,
		description: "convert lists to sets",
		hard:        true,
		apply: func(in map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
			var cmd struct {
				Name        string
				User        string
				Source      string
				Environment map[string]string `dynamodbav:",omitempty"`
				ACL         []string          `dynamodbav:",stringset,omitempty"`
				Admins      []string          `dynamodbav:",stringset,omitempty"`
				Description string            `dynamodbav:",omitempty"`
			}
			if err := dynamodbattribute.UnmarshalMap(in, &cmd); err != nil {
				return nil, err
			}

			return dynamodbattribute.MarshalMap(cmd)
		},
	},
}

// Migrations contains a map of migrations keyed by the version number
type Migrations map[int]migration

// Apply all migrations on copy of item until item version is equal to target.
// Attempting to apply a hard migration when allowHard is false will result in
// an error.
func (ms Migrations) Apply(target int, allowHard bool, item map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	initialVer := itemVersion(item)
	cur := item
	for nextVer := initialVer + 1; nextVer <= target; nextVer++ {
		m, ok := ms[nextVer]
		if !ok {
			continue
		}

		if !allowHard && m.hard {
			return cur, errors.Errorf("hard migrations disabled version: %d - %s", nextVer, m.description)
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

// HardRequired returns true if a hard migration will be required to migrate
// to the target schema version.
func (ms Migrations) HardRequired(current, target int) bool {
	for v := current + 1; v <= target; v++ {
		m, ok := ms[v]
		if ok && m.hard {
			return true
		}
	}
	return false
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

// migrateTable to target version by calling migrations.Apply on each item.
// Upon completion schemaAttr tag will reflect the target schema version unless
// the client region is 'local'.
func migrateTable(client *dynamodb.DynamoDB, table string, target int) error {
	var items []map[string]*dynamodb.AttributeValue
	var ops []*dynamodb.WriteRequest

	// scan all items in table.
	// NOTE: setting a pagesize may be required if throttling occurs
	err := client.ScanPages(&dynamodb.ScanInput{
		TableName:      aws.String(table),
		ConsistentRead: aws.Bool(true),
	}, func(p *dynamodb.ScanOutput, lastPage bool) bool {
		items = append(items, p.Items...)
		return !lastPage
	})
	if err != nil {
		return err
	}
	log.Info("migrating", log.Fields{"items": strconv.Itoa(len(items))})
	for _, item := range items {
		if itemVersion(item) == target {
			// skip items at target version avoiding unnescesary PutRequests.
			continue
		}
		mItem, err := migrations.Apply(target, true, item)
		if err != nil {
			return err
		}
		// append migrated item to ops slice
		ops = append(ops, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: mItem,
			},
		})
	}

	batches := int(math.Ceil(float64(len(ops)) / maxBatchSize))
	boff := backoff.NewExponentialBackOff()
	for i := 0; i < batches; i++ {
		start, end := i*maxBatchSize, (i+1)*maxBatchSize
		if end > len(ops) {
			end = len(ops)
		}
		curOps := ops[start:end]
		out, err := client.BatchWriteItem(&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				table: curOps,
			},
		})
		if err != nil {
			return err
		}

		if len(out.UnprocessedItems) != 0 {
			fmt.Printf("%#v\n", out.UnprocessedItems)
			// fail since we don't retry any UnprocessedItems.
			return errors.New("failed to process items")
		}
		time.Sleep(boff.NextBackOff())
	}

	if aws.StringValue(client.Config.Region) == "local" {
		log.Info("skipping unsupported dynamodb-local operation", log.Fields{"operation": "setTableVersion"})
		return nil
	}
	return setTableVersion(client, table, target)
}
