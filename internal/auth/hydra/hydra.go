package hydra

import (
	"context"
	"errors"

	client "github.com/ory/hydra-client-go/v2"
	"github.com/saygik/go-userinfo/internal/entity"
)

type Hydra struct {
	client *client.APIClient
	scopes []entity.IDPScope
}

func New(client *client.APIClient, scopes []entity.IDPScope) *Hydra {
	return &Hydra{
		client: client,
		scopes: scopes,
	}
}
func (hdr Hydra) GetScopes() []entity.IDPScope {
	return hdr.scopes
}
func (hdr Hydra) CheckHydra() bool {
	_, _, err := hdr.client.WellknownAPI.DiscoverJsonWebKeys(context.Background()).Execute()

	return err == nil
}

func (hdr Hydra) LogoutURL() string {
	resp, _, _ := hdr.client.OidcAPI.DiscoverOidcConfiguration(context.Background()).Execute()
	return *resp.EndSessionEndpoint
}
func (hdr Hydra) AcceptOAuth2LogoutRequest(logoutChallenge string) (string, error) {
	//	   r, err := apiClient.OidcAPI.RevokeOidcSession(context.Background()).Execute()
	//	_, resp := hdr.client.OAuth2API.RevokeOAuth2ConsentSessions(context.Background()).Execute()
	resp, _, err := hdr.client.OAuth2API.AcceptOAuth2LogoutRequest(context.Background()).LogoutChallenge(logoutChallenge).Execute()

	if err != nil {
		return "", errors.New("cannot accept logout request`")
	}

	//_, _, resp := hdr.client.OAuth2API.AcceptOAuth2LogoutRequest(context.Background()).Execute()
	_ = resp
	return resp.RedirectTo, nil
}

func (hdr Hydra) GetOAuth2LoginRequest(loginChallenge string) (*entity.OAuth2LoginRequest, error) {
	resp, _, err := hdr.client.OAuth2API.GetOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).Execute()
	if err != nil {
		return nil, err
	}
	// resp1, _, err := hdr.client.OAuth2API.ListOAuth2ConsentSessions(context.Background()).Subject(resp.Subject).PageSize(int64(789)).PageToken("pageToken_example").LoginSessionId(*resp.SessionId).Execute()
	// _ = resp1
	return &entity.OAuth2LoginRequest{
		Challenge:                    resp.Challenge,
		RequestUrl:                   resp.RequestUrl,
		RequestedAccessTokenAudience: resp.RequestedAccessTokenAudience,
		RequestedScope:               resp.RequestedScope,
		SessionId:                    resp.SessionId,
		Skip:                         resp.Skip,
		Subject:                      resp.Subject,
		AdditionalProperties:         resp.AdditionalProperties,
	}, nil
}

func (hdr Hydra) AcceptNewOAuth2LoginRequest(loginChallenge string, subject string, rememberMe bool) (string, error) {
	acceptOAuth2LoginRequest := *client.NewAcceptOAuth2LoginRequest(subject) // AcceptOAuth2LoginRequest |  (optional)
	acceptOAuth2LoginRequest.SetRemember(rememberMe)
	acceptOAuth2LoginRequest.SetRememberFor(2678400)

	resp, _, err := hdr.client.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).AcceptOAuth2LoginRequest(acceptOAuth2LoginRequest).Execute()
	if err != nil {
		return "", errors.New("cannot accept new login request`")
	}
	return resp.RedirectTo, nil
}

func (hdr Hydra) AcceptOAuth2LoginRequest(loginChallenge string, subject string) (string, error) {
	acceptOAuth2LoginRequest := *client.NewAcceptOAuth2LoginRequest(subject)
	resp2, _, err := hdr.client.OAuth2API.AcceptOAuth2LoginRequest(context.Background()).LoginChallenge(loginChallenge).AcceptOAuth2LoginRequest(acceptOAuth2LoginRequest).Execute()
	if err != nil {
		return "", errors.New("cannot Accept Login Request")
	}
	return resp2.RedirectTo, nil
}

func (hdr Hydra) GetOAuth2ConsentRequest(consentChallenge string) (*entity.OAuth2ConsentRequest, error) {
	resp, _, err := hdr.client.OAuth2API.GetOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentChallenge).Execute()
	if err != nil {
		return nil, errors.New("cannot Accept Consent Request`")
	}
	skipConsent := *resp.Skip || *resp.Client.SkipConsent
	return &entity.OAuth2ConsentRequest{
		Challenge:                    resp.Challenge,
		RequestUrl:                   resp.RequestUrl,
		RequestedAccessTokenAudience: resp.RequestedAccessTokenAudience,
		RequestedScope:               resp.RequestedScope,
		Skip:                         &skipConsent,
		Subject:                      resp.Subject,
		AdditionalProperties:         resp.AdditionalProperties,
	}, nil
}

func (hdr Hydra) GetOAuth2Client(clientId string) (*entity.OAuth2Client, error) {

	client, _, err := hdr.client.OAuth2API.GetOAuth2Client(context.Background(), clientId).Execute()
	if err != nil {
		return nil, errors.New("не найдена система-клиент в oAuth2")
	}
	_ = client
	resp := entity.OAuth2Client{
		ClientId:   client.ClientId,
		ClientName: client.ClientName,
	}
	return &resp, nil
}

func (hdr Hydra) AcceptOAuth2ConsentRequest(consentRequest *entity.OAuth2ConsentRequest, user map[string]interface{}) (string, error) {
	acceptOAuth2ConsentRequest := *client.NewAcceptOAuth2ConsentRequest()
	newAcceptOAuth2ConsentRequestSession := *client.NewAcceptOAuth2ConsentRequestSession()
	acceptOAuth2ConsentRequest.SetGrantAccessTokenAudience(consentRequest.RequestedAccessTokenAudience)
	acceptOAuth2ConsentRequest.SetGrantScope(consentRequest.RequestedScope)

	claims := map[string]interface{}{}

	for _, scope := range consentRequest.RequestedScope {
		if scope == "email" {
			claims["email"] = user["mail"]
		}
		if scope == "profile" {
			claims["family_name"] = user["sn"]
			claims["given_name"] = user["givenName"]
			claims["name"] = user["displayName"]
			claims["company"] = user["company"]
			claims["department"] = user["department"]
			claims["title"] = user["title"]
			claims["nickname"] = user["userPrincipalName"]
			claims["phone_number"] = user["telephoneNumber"]
		}
		if scope == "role" {
			claims["roles"] = user["roles"]
		}
	}
	newAcceptOAuth2ConsentRequestSession.IdToken = claims
	acceptOAuth2ConsentRequest.SetSession(newAcceptOAuth2ConsentRequestSession)
	resp2, _, err := hdr.client.OAuth2API.AcceptOAuth2ConsentRequest(context.Background()).ConsentChallenge(consentRequest.Challenge).AcceptOAuth2ConsentRequest(acceptOAuth2ConsentRequest).Execute()
	if err != nil {
		return "", errors.New("cannot Accept Consent Request")
	}
	return resp2.RedirectTo, nil
}

func (hdr Hydra) IntrospectOAuth2Token(token string) (*entity.IntrospectedOAuth2Token, error) {

	resp, _, err := hdr.client.OAuth2API.IntrospectOAuth2Token(context.Background()).Token(token).Scope("").Execute()
	_ = resp

	if err != nil {
		return nil, errors.New("cannot introspect token`")
	}
	introspectedOAuth2Token := entity.IntrospectedOAuth2Token(*resp)

	return &introspectedOAuth2Token, nil
}

// func (hdr Hydra) RevokeOidcSession(subject string) error {
// 	_, err := hdr.client.OidcAPI.RevokeOidcSession(context.Background()).Execute()
// 	if err != nil {
// 		return err
// 	}
// 	_, err = hdr.client.OAuth2API.RevokeOAuth2ConsentSessions(context.Background()).Subject(subject).Client("client103").All(true).Execute()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
//}
