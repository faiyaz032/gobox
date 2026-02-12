package docker

import (
	"github.com/docker/docker/client"
)

type Svc struct {
	client *client.Client
}

func NewSvc() (*Svc, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Svc{client: cli}, nil
}

func (s *Svc) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}
