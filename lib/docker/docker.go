package docker

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/gliderlabs/comlab/pkg/log"
)

func init() {
	com.Register("docker", &Component{},
		com.Option("name", "", "srv record for docker host discovery"),
		com.Option("version", client.DefaultVersion, "Docker client API version"))
}

var clients []client.APIClient

type Component struct{}

// Client returns a random docker APIClient
func Client() client.APIClient {
	return clients[rand.Intn(len(clients))]
}

// AppPreStart sets up docker client pool
func (c *Component) AppPreStart() error {
	rand.Seed(time.Now().UTC().UnixNano())

	if name := com.GetString("name"); name != "" {
		cli, err := discoverHosts(name)
		if err != nil {
			return err
		}
		clients = append(clients, cli...)
		return nil
	}
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	log.Info("using env client")
	clients = append(clients, cli)
	return nil
}

// discoverHosts from srv lookup on name.
func discoverHosts(name string) ([]client.APIClient, error) {
	_, addrs, err := net.LookupSRV("", "", name)
	if err != nil {
		return nil, err
	}
	var clients []client.APIClient
	for _, addr := range addrs {
		target := strings.TrimSuffix(addr.Target, ".")
		host := fmt.Sprintf("tcp://%s:%v", target, addr.Port)
		apiversion := com.GetString("version")
		c, err := client.NewClient(host, apiversion, nil, nil)
		if err != nil {
			log.Info("failed to create docker client for: "+addr.Target, err)
			continue
		}
		// retrieve version to test connectivity
		version, err := c.ServerVersion(context.Background())
		if err != nil {
			log.Info("failed to retrieve docker version for: "+addr.Target, err)
			continue
		}
		log.Info(log.Fields{
			"host":        addr.Target,
			"api.version": version.APIVersion,
			"version":     version.Version})
		clients = append(clients, c)
	}
	if len(clients) == 0 {
		return nil, fmt.Errorf("no clients")
	}
	return clients, nil
}
