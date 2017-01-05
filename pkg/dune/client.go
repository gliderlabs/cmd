package dune

import (
	"fmt"
	"net"
	"strings"

	"github.com/docker/docker/client"
)

type Client struct {
	*client.Client
	host string
}

func NewClient(name string) (*Client, error) {
	_, addrs, err := net.LookupSRV("", "", name)
	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		return nil, fmt.Errorf("err: no addrs found")
	}

	c := &Client{}
	c.host = strings.TrimSuffix(addrs[0].Target, ".")
	c.Client, err = client.NewClient(fmt.Sprintf("tcp://%s:%v", c.host, addrs[0].Port), client.DefaultVersion, nil, nil)
	return c, err
}

func NewEnvClient() (*Client, error) {
	c := &Client{}
	var err error
	c.Client, err = client.NewEnvClient()
	return c, err
}

// Host returns docker host
func (c *Client) Host() string {
	return c.host
}
