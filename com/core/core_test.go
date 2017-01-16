package core

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/progrium/cmd/pkg/dune"
	mock_client "github.com/progrium/cmd/pkg/mock/docker/docker/client"
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

func TestCmdPull(t *testing.T) {
	cmd := &Command{
		Source: "progrium/cmd",
		Name:   "alpine",
		User:   "nobody",
	}

	t.Run("ExcedeLimit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mock_client.NewMockAPIClient(ctrl)
		cmd.docker = &dune.Client{APIClient: client}
		pullRes := ioutil.NopCloser(strings.NewReader(""))
		client.EXPECT().
			ImagePull(gomock.Any(), cmd.Source, types.ImagePullOptions{}).
			Return(pullRes, nil)

		// Return an image size 1 byte larger than Plans[DefaultPlan].ImageSize
		client.EXPECT().
			ImageInspectWithRaw(gomock.Any(), cmd.Source).
			Return(types.ImageInspect{Size: Plans[DefaultPlan].ImageSize + 1}, []byte{}, nil)
		client.EXPECT().
			ImageRemove(gomock.Any(), cmd.Source, types.ImageRemoveOptions{}).
			Return([]types.ImageDelete{}, nil)

		assert.Error(t, cmd.Pull())
	})

	t.Run("WithinLimit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mock_client.NewMockAPIClient(ctrl)
		cmd.docker = &dune.Client{APIClient: client}
		pullRes := ioutil.NopCloser(strings.NewReader(""))
		client.EXPECT().
			ImagePull(gomock.Any(), cmd.Source, types.ImagePullOptions{}).
			Return(pullRes, nil)
		client.EXPECT().
			ImageInspectWithRaw(gomock.Any(), cmd.Source).
			Return(types.ImageInspect{Size: Plans[DefaultPlan].ImageSize}, []byte{}, nil)
		client.EXPECT().
			ImageTag(gomock.Any(), cmd.Source, cmd.image()).
			Return(nil)

		assert.NoError(t, cmd.Pull())
	})

}
