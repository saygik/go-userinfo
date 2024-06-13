package app

import (
	"context"
	"fmt"

	client "github.com/ory/hydra-client-go/v2"
)

func (a *App) newHydraClient(url string) (*client.APIClient, error) {

	configuration := client.NewConfiguration()
	configuration.Servers = []client.ServerConfiguration{
		{
			URL: url, // Admin API URL
		},
	}
	cl := client.NewAPIClient(configuration)

	_, r, err := cl.WellknownAPI.DiscoverJsonWebKeys(context.Background()).Execute()
	if err != nil {
		return nil, fmt.Errorf("oAuth2 Hydra - unable to get Wellknown API: %v\n", r)
	}
	return cl, nil
}
