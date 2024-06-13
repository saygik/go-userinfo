package hydra

import (
	"context"

	client "github.com/ory/hydra-client-go/v2"
)

type Hydra struct {
	client *client.APIClient
}

func New(client *client.APIClient) *Hydra {
	return &Hydra{
		client: client,
	}
}

func (hdr Hydra) CheckHydra() bool {
	_, _, err := hdr.client.WellknownAPI.DiscoverJsonWebKeys(context.Background()).Execute()
	return err == nil
}
