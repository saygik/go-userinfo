package oauth2

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc"
	"github.com/saygik/go-userinfo/internal/entity"
	"golang.org/x/oauth2"
)

type OAuth2 struct {
	oAuth2Config oauth2.Config
	oidcVerifier *oidc.IDTokenVerifier
	oidcProvider *oidc.Provider
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func New(oAuth2Config oauth2.Config, oidcProvider *oidc.Provider) *OAuth2 {
	oidcConfig := &oidc.Config{
		ClientID: oAuth2Config.ClientID,
	}

	verifier := oidcProvider.Verifier(oidcConfig)
	return &OAuth2{
		oAuth2Config: oAuth2Config,
		oidcVerifier: verifier,
		oidcProvider: oidcProvider,
	}

}
func (o OAuth2) GetRedirectUrl() string {
	return o.oAuth2Config.RedirectURL
}

func (o OAuth2) AuthCodeURL(state string) string {
	return o.oAuth2Config.AuthCodeURL(state)
}

func (o OAuth2) Refresh(rtoken string) (entity.Token, error) {
	token := entity.Token{}
	t := &oauth2.Token{
		RefreshToken: rtoken,
	}
	tokenSource := o.oAuth2Config.TokenSource(context.Background(), t)
	newToken, err := tokenSource.Token()

	if err != nil {
		return token, err
	}
	if newToken.RefreshToken == rtoken {
		return token, errors.New("token does not refreshed")
	}
	token.AccessToken = newToken.AccessToken
	token.RefreshToken = newToken.RefreshToken
	token.IdToken = newToken.Extra("id_token").(string)
	token.Expiry = newToken.Expiry

	return token, nil
}

func (o OAuth2) Exchange(code string) (*entity.Token, *entity.UserInfo, error) {
	accessToken, err := o.oAuth2Config.Exchange(context.Background(), code)
	token := entity.Token{}

	if err != nil {
		return nil, nil, err
	}

	userInfo, err := o.oidcProvider.UserInfo(context.Background(), oauth2.StaticTokenSource(accessToken))
	_ = userInfo
	if err != nil {
		return nil, nil, err
	}
	user := entity.UserInfo{}
	err = userInfo.Claims(&user)
	if err != nil {
		return nil, nil, err
	}

	rawIDToken, ok := accessToken.Extra("id_token").(string)
	if !ok {
		return nil, nil, errors.New("невозможно определить id token")
	}

	// 	idToken, err := o.oidcVerifier.Verify(context.Background(), rawIDToken)
	// if err != nil {
	// 	return token, errors.New("невозможно определить id token")
	// }
	token.AccessToken = accessToken.AccessToken
	token.RefreshToken = accessToken.RefreshToken
	token.IdToken = rawIDToken
	token.Expiry = accessToken.Expiry

	return &token, &user, nil

}
