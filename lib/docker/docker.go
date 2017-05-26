package docker

import (
	"log"
	"strings"

	"github.com/docker/docker/client"
	"github.com/gliderlabs/comlab/pkg/com"
	"github.com/pborman/uuid"
)

// Client returns a sandbox docker APIClient
func Client() client.APIClient {
	c, err := getClient()
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func getClient() (client.APIClient, error) {
	if com.GetString("name") == "" {
		return client.NewEnvClient()
	}
	var (
		host    = "tcp://" + com.GetString("listen")
		version = com.GetString("version")
		headers = map[string]string{
			SessionHeaderKey: strings.SplitN(uuid.New(), "-", 2)[0],
		}
	)
	return client.NewClient(host, version, nil, headers)
}
