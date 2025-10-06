package oauth2authentik

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/saygik/go-userinfo/internal/entity"
	"golang.org/x/oauth2"
)

type OAuth2 struct {
	oAuth2Config oauth2.Config
	oidcProvider *oidc.Provider
	logoutUrl    string
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func New(oAuth2Config oauth2.Config, oidcProvider *oidc.Provider, logoutUrl string) *OAuth2 {

	return &OAuth2{
		oAuth2Config: oAuth2Config,
		oidcProvider: oidcProvider,
		logoutUrl:    logoutUrl,
	}
}

func (o OAuth2) GetRedirectUrl() string {
	return o.oAuth2Config.RedirectURL
}

func (o OAuth2) AuthCodeURL(state string) string {
	return o.oAuth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}
func (o OAuth2) LogOutURL() string {
	return o.logoutUrl
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
	accessToken, err := o.oAuth2Config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("access_type", "offline"))
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

	// rawIDToken, ok := accessToken.Extra("id_token").(string)
	// if !ok {
	// 	return nil, nil, errors.New("невозможно определить id token")
	// }

	// 	idToken, err := o.oidcVerifier.Verify(context.Background(), rawIDToken)
	// if err != nil {
	// 	return token, errors.New("невозможно определить id token")
	// }
	token.AccessToken = accessToken.AccessToken
	token.RefreshToken = accessToken.RefreshToken
	//	token.IdToken = rawIDToken
	token.Expiry = accessToken.Expiry

	return &token, &user, nil

}

type jwtStdClaims struct {
	Exp   int64  `json:"exp"`
	Sub   string `json:"sub"`   // standard user id
	Name  string `json:"name"`  // if your provider uses this
	Email string `json:"email"` // standard user Email
}

func (o OAuth2) ExchangeRefreshToAccessToken(refreshToken string) (*entity.Token, error) {
	ctx := context.Background()
	ts := o.oAuth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})

	tok, err := ts.Token()
	if err != nil {
		return nil, err
	}

	return &entity.Token{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
	}, nil

}

func (o OAuth2) IntrospectOAuth2Token(querytoken string) (*entity.UserInfo, error) {
	//	token := &oauth2.Token{AccessToken: querytoken}
	parts := strings.Split(querytoken, ".")

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var claims jwtStdClaims

	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	expiry := time.Unix(claims.Exp, 0)
	if time.Now().After(expiry) {
		return nil, errors.New("token expired")
	}

	// userInfo, err := o.oidcProvider.UserInfo(context.Background(), oauth2.StaticTokenSource(token))
	// _ = userInfo
	// if err != nil {
	// 	return nil, err
	// }
	user := entity.UserInfo{}
	//err = userInfo.Claims(&user)
	// if err != nil {
	// 	return nil, err
	// }
	user.Sub = claims.Sub
	user.Name = claims.Name
	user.Email = claims.Email
	return &user, err
}
