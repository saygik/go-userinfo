package app

import (
	"context"
	"log"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type OAuth2Client struct {
	oAuth2Config oauth2.Config
	oidcProvider *oidc.Provider
	logoutUrl    string
}

// type OAuth2ClientAuthentik struct {
// 	oAuth2Config oauth2.Config
// }

func (a *App) newOAuth2Client(url string, clientID string, clientSecret string, redirectUrl string, scopes []string, logoutUrl string) (*OAuth2Client, error) {

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, url)
	if err != nil {
		log.Fatal(err)
	}

	oAuthConf2 := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectUrl,
		Scopes:       scopes,
	}
	return &OAuth2Client{
		oAuth2Config: oAuthConf2,
		oidcProvider: provider,
		logoutUrl:    logoutUrl,
	}, nil
}
