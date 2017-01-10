package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasAccess(t *testing.T) {

	var testCases = []struct {
		Cmd      *Command
		User     string
		Expected bool
	}{
		{
			Cmd:      &Command{},
			User:     "nobody",
			Expected: false,
		},
		{
			Cmd:      &Command{User: "nobody"},
			User:     "nobody",
			Expected: true,
		},
		{
			Cmd:      &Command{User: "somebody"},
			User:     "nobody",
			Expected: false,
		},
		{
			Cmd:      &Command{Admins: []string{"nobody", "somebody"}},
			User:     "nobody",
			Expected: true,
		},
		{
			Cmd:      &Command{Admins: []string{"somebody"}},
			User:     "nobody",
			Expected: false,
		},
		{
			Cmd:      &Command{ACL: []string{"nobody", "somebody"}},
			User:     "nobody",
			Expected: true,
		},
		{
			Cmd:      &Command{ACL: []string{"somebody"}},
			User:     "nobody",
			Expected: false,
		},
		{
			Cmd:      &Command{ACL: []string{"*"}}, // Public ACL
			User:     "nobody",
			Expected: false,
		},
	}

	for _, test := range testCases {
		actual := test.Cmd.HasAccess(test.User)
		assert.Equal(t, test.Expected, actual)
	}
}
