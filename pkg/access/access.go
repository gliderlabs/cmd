package access

import "github.com/gliderlabs/comlab/pkg/log"

// Allows the component package to hook in and
// provide the default client.
var DefaultClientFactory func() Client

// Since DefaultClient is a function, we memoize
// the client after first return to make it a
// singleton.
var defaultClient Client

func Check(name string) bool {
	if defaultClient != nil {
		return defaultClient.Check(name)
	}
	if DefaultClientFactory != nil {
		defaultClient = DefaultClientFactory()
	} else {
		log.Info("DefaultClientFactory unregistered: using denyClient")
		defaultClient = &denyClient{}
	}
	return defaultClient.Check(name)
}

type Client interface {
	Check(name string) bool
}

type denyClient struct{}

func (denyClient) Check(name string) bool {

	return false
}
