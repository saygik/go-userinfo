package ad

import (
	adClient "github.com/saygik/go-ad-client"
)

type Repository struct {
	ads map[string]*adClient.ADClient
}

func New(adclients map[string]*adClient.ADClient) *Repository {
	return &Repository{
		ads: adclients,
	}
}
