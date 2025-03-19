package ad

import (
	"github.com/saygik/go-userinfo/internal/config"
	adClient "github.com/saygik/go-userinfo/pkg/ad-client"
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
