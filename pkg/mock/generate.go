package mock

//go:generate mockgen -package client -destination docker/docker/client/api.go github.com/docker/docker/client APIClient
