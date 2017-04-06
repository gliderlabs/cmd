package core

import (
	"context"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	mock_client "github.com/gliderlabs/cmd/lib/mock/docker/docker/client"
	"github.com/stretchr/testify/assert"

	"github.com/gliderlabs/cmd/app/billing"
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
func TestParseSource(t *testing.T) {
	var testCases = []struct {
		Source      []byte
		ExpectImage string
		ExpectPkgs  []string
		ExpectBody  []byte
		ExpectErr   bool
	}{
		{
			Source:    []byte("#!cmd alpine"),
			ExpectErr: true,
		},
		{
			Source:    []byte("\n#!cmd alpine\n"), // #!cmd must be on first line
			ExpectErr: true,
		},
		{
			Source:    []byte("#!/usr/bin/bash\n"),
			ExpectErr: true,
		},
		{
			Source:      []byte("#!cmd alpine\n"),
			ExpectImage: "alpine",
			ExpectBody:  []byte{},
		},
		{
			Source:      []byte("#!cmd alpine bash\n"),
			ExpectImage: "alpine",
			ExpectPkgs:  []string{"bash"},
			ExpectBody:  []byte{},
		},
		{
			Source:      []byte("#!cmd alpine bash\n#!/usr/bin/bash"),
			ExpectImage: "alpine",
			ExpectPkgs:  []string{"bash"},
			ExpectBody:  []byte("#!/usr/bin/bash"),
		},
		{
			Source:      []byte("#!cmd alpine bash\n\n#!/usr/bin/bash\n"),
			ExpectImage: "alpine",
			ExpectPkgs:  []string{"bash"},
			ExpectBody:  []byte("\n#!/usr/bin/bash\n"),
		},
		{
			Source:      []byte("#!cmd alpine bash\n#!/usr/bin/bash\necho"),
			ExpectImage: "alpine",
			ExpectPkgs:  []string{"bash"},
			ExpectBody:  []byte("#!/usr/bin/bash\necho"),
		},
		{
			Source:      []byte("#!cmd nate/pandashells bash go\n"),
			ExpectImage: "nate/pandashells",
			ExpectPkgs:  []string{"bash", "go"},
			ExpectBody:  []byte{},
		},
		{
			Source: []byte(`#!cmd nate/pandashells
#!/bin/bash
/usr/local/bin/p.example_data "$@"`),
			ExpectImage: "nate/pandashells",
			ExpectBody: []byte(`#!/bin/bash
/usr/local/bin/p.example_data "$@"`),
		},
	}

	for i, test := range testCases {
		t.Run("Case:"+strconv.Itoa(i+1), func(t *testing.T) {
			img, pkgs, body, err := parseSource(test.Source)
			if test.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.ExpectImage, img)
			assert.Equal(t, test.ExpectPkgs, pkgs)
			assert.Equal(t, test.ExpectBody, body)
		})
	}
}

func TestMakeBuildCtx(t *testing.T) {
	var testCases = []struct {
		Image     string
		Pkgs      []string
		Body      []byte
		ExpectCtx map[string][]byte
		ExpectErr bool
	}{
		{
			Image:     "alpine", // Missing body
			ExpectErr: true,
			ExpectCtx: map[string][]byte{},
		},
		{
			Image:     "nate/pandashells", // Unsupported image
			Body:      []byte("#!/usr/bin/bash\n"),
			ExpectErr: true,
			ExpectCtx: map[string][]byte{},
		},
		{
			Image: "alpine",
			Body:  []byte("#!/usr/bin/bash\n"),
			ExpectCtx: map[string][]byte{
				"Dockerfile": []byte("FROM alpine\nWORKDIR /cmd\nENTRYPOINT [\"/usr/bin/bash\"]\n"),
			},
		},
		{
			Image: "alpine",
			Pkgs:  []string{"bash"},
			Body:  []byte("#!/usr/bin/bash\n"),
			ExpectCtx: map[string][]byte{
				"Dockerfile": []byte(`FROM alpine
RUN apk --no-cache add bash
WORKDIR /cmd
ENTRYPOINT ["/usr/bin/bash"]
`),
			},
		},
		{
			Image: "alpine",
			Body:  []byte("#!/usr/bin/bash\necho"),
			ExpectCtx: map[string][]byte{
				"Dockerfile": []byte(`FROM alpine
COPY ./entrypoint ./bin/entrypoint
WORKDIR /cmd
ENTRYPOINT ["/bin/entrypoint"]
`),
				"entrypoint": []byte("#!/usr/bin/bash\necho"),
			},
		},
	}

	for _, test := range testCases {
		ctx, err := getBuildCtx(test.Image, test.Pkgs, test.Body)
		if test.ExpectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		assert.Equal(t, test.ExpectCtx, ctx)

	}
}

func TestCmdPull(t *testing.T) {
	cmd := &Command{
		Source: "gliderlabs/cmd",
		Name:   "alpine",
		User:   "nobody",
	}

	t.Run("ExceedLimit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mock_client.NewMockAPIClient(ctrl)
		cmd.docker = client
		pullRes := ioutil.NopCloser(strings.NewReader(""))
		client.EXPECT().
			ImagePull(gomock.Any(), cmd.Source, types.ImagePullOptions{}).
			Return(pullRes, nil)

		// Return an image size 1 byte larger than Plans[DefaultPlan].ImageSize
		client.EXPECT().
			ImageInspectWithRaw(gomock.Any(), cmd.Source).
			Return(types.ImageInspect{Size: billing.Plans[billing.DefaultPlan].ImageSize + 1}, []byte{}, nil)
		client.EXPECT().
			ImageRemove(gomock.Any(), cmd.Source, types.ImageRemoveOptions{}).
			Return([]types.ImageDeleteResponseItem{}, nil)

		assert.Error(t, cmd.Pull(context.Background()))
	})

	t.Run("WithinLimit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mock_client.NewMockAPIClient(ctrl)
		cmd.docker = client
		pullRes := ioutil.NopCloser(strings.NewReader(""))
		client.EXPECT().
			ImagePull(gomock.Any(), cmd.Source, types.ImagePullOptions{}).
			Return(pullRes, nil)
		client.EXPECT().
			ImageInspectWithRaw(gomock.Any(), cmd.Source).
			Return(types.ImageInspect{Size: billing.Plans[billing.DefaultPlan].ImageSize}, []byte{}, nil)
		client.EXPECT().
			ImageTag(gomock.Any(), cmd.Source, cmd.image()).
			Return(nil)

		assert.NoError(t, cmd.Pull(context.Background()))
	})

}
