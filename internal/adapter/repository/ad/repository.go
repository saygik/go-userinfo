package ad

import (
	adClient "github.com/saygik/go-ad-client"
	"github.com/saygik/go-userinfo/internal/config"
)

type Repository struct {
	ads       map[string]*adClient.ADClient
	adconfigs map[string]*config.ADConfig
}

func New(adclients map[string]*adClient.ADClient, adconfigs map[string]*config.ADConfig) *Repository {
	return &Repository{
		ads:       adclients,
		adconfigs: adconfigs,
	}
}
