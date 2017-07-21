package dockerbox

import (
	"fmt"
	"net"

	"github.com/docker/docker/client"
	"github.com/gliderlabs/comlab/pkg/com"
)

const (
	APIVersion = "1.27" // TODO: check dockerbox
)

func init() {
	com.Register("dockerbox", &Component{},
		com.Option("hostname", "", "hostname used to get backend A records"),
	)
}

type Component struct{}

type Client struct {
	client.APIClient
	Host string
}

func GetBackend() (*Client, error) {
	if com.GetString("hostname") == "" {
		c, err := client.NewEnvClient()
		return &Client{c, "local"}, err
	}
	addrs, err := net.LookupHost(com.GetString("hostname"))
	if err != nil {
		return nil, err
	}
	c, err := client.NewClient(fmt.Sprintf("tcp://%s:2375", addrs[0]), APIVersion, nil, nil)
	return &Client{c, addrs[0]}, err
}
