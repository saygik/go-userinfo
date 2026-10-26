package oauth2authentik

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/saygik/go-userinfo/internal/entity"
	"golang.org/x/oauth2"
)

// OAuth2 - клиент Authentik: вход по коду, обновление и отзыв токенов,
// проверка access-токена входящих запросов.
type OAuth2 struct {
	oAuth2Config  oauth2.Config
	oidcProvider  *oidc.Provider
	verifier      *oidc.IDTokenVerifier
	mainIssuer    string
	verifiers     *issuerVerifierCache
	logoutUrl     string
	revocationURL string
}

// issuerVerifierCache хранит verifier'ы других приложений того же Authentik
// (например, Grafana). В Authentik у каждого приложения свой issuer и свой JWKS,
// поэтому одним verifier'ом токены разных приложений проверить нельзя.
type issuerVerifierCache struct {
	mu sync.RWMutex
	m  map[string]*oidc.IDTokenVerifier
}

const authentikAppsPath = "/application/o/"

var issuerSlugRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]*$`)

func New(oAuth2Config oauth2.Config, oidcProvider *oidc.Provider, logoutUrl string) *OAuth2 {

	// Подпись токена проверяется по JWKS провайдера. Проверку срока действия
	// выполняем сами: для запросов от доверенных API она отключена.
	verifier := oidcProvider.Verifier(&oidc.Config{
		ClientID:          oAuth2Config.ClientID,
		SkipClientIDCheck: true,
		SkipExpiryCheck:   true,
	})

	var providerClaims struct {
		Issuer             string `json:"issuer"`
		RevocationEndpoint string `json:"revocation_endpoint"`
	}
	_ = oidcProvider.Claims(&providerClaims)

	return &OAuth2{
		oAuth2Config:  oAuth2Config,
		oidcProvider:  oidcProvider,
		verifier:      verifier,
		mainIssuer:    providerClaims.Issuer,
		verifiers:     &issuerVerifierCache{m: map[string]*oidc.IDTokenVerifier{providerClaims.Issuer: verifier}},
		logoutUrl:     logoutUrl,
		revocationURL: providerClaims.RevocationEndpoint,
	}
}

func (o OAuth2) GetRedirectUrl() string {
	return o.oAuth2Config.RedirectURL
}

func (o OAuth2) LogOutURL() string {
	return o.logoutUrl
}

// AuthCodeURL собирает ссылку на вход в Authentik.
// codeChallenge - PKCE-преобразование (S256) от code_verifier.
func (o OAuth2) AuthCodeURL(state, codeChallenge string) string {
	return o.oAuth2Config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
}

// CodeChallenge считает PKCE S256 от code_verifier.
func (o OAuth2) CodeChallenge(codeVerifier string) string {
	sum := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// Exchange меняет код авторизации на токены и забирает данные пользователя.
func (o OAuth2) Exchange(code, codeVerifier string) (*entity.Token, *entity.UserInfo, error) {
	ctx := context.Background()

	opts := make([]oauth2.AuthCodeOption, 0, 1)
	if codeVerifier != "" {
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	}

	accessToken, err := o.oAuth2Config.Exchange(ctx, code, opts...)
	if err != nil {
		return nil, nil, err
	}

	userInfo, err := o.oidcProvider.UserInfo(ctx, oauth2.StaticTokenSource(accessToken))
	if err != nil {
		return nil, nil, err
	}
	user := entity.UserInfo{}
	if err := userInfo.Claims(&user); err != nil {
		return nil, nil, err
	}

	token := entity.Token{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: accessToken.RefreshToken,
		Expiry:       accessToken.Expiry,
	}
	if idToken, ok := accessToken.Extra("id_token").(string); ok {
		token.IdToken = idToken
	}

	return &token, &user, nil
}

// Refresh обновляет токены по refresh-токену.
// Если провайдер не ротирует refresh-токен, возвращается прежний.
func (o OAuth2) Refresh(refreshToken string) (entity.Token, error) {
	if refreshToken == "" {
		return entity.Token{}, fmt.Errorf("%w: empty refresh token", entity.ErrRefreshTokenRejected)
	}

	ctx := context.Background()
	newToken, err := o.oAuth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken}).Token()
	if err != nil {
		return entity.Token{}, classifyRefreshError(err)
	}

	token := entity.Token{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		Expiry:       newToken.Expiry,
	}
	if token.RefreshToken == "" {
		token.RefreshToken = refreshToken
	}
	if idToken, ok := newToken.Extra("id_token").(string); ok {
		token.IdToken = idToken
	}

	return token, nil
}

func (o OAuth2) ExchangeRefreshToAccessToken(refreshToken string) (*entity.Token, error) {
	token, err := o.Refresh(refreshToken)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// classifyRefreshError отличает «токен отвергнут» (нужен повторный вход)
// от временных проблем провайдера (сеть, 5xx): во втором случае сессию
// пользователя завершать нельзя.
func classifyRefreshError(err error) error {
	var retrieveErr *oauth2.RetrieveError
	if errors.As(err, &retrieveErr) {
		switch retrieveErr.ErrorCode {
		case "invalid_grant", "invalid_token":
			return fmt.Errorf("%w: %s", entity.ErrRefreshTokenRejected, err)
		}
		if retrieveErr.Response != nil {
			code := retrieveErr.Response.StatusCode
			if code == http.StatusBadRequest || code == http.StatusUnauthorized {
				return fmt.Errorf("%w: %s", entity.ErrRefreshTokenRejected, err)
			}
		}
	}
	return err
}

// IntrospectOAuth2Token проверяет подпись, издателя и формат access-токена.
// Срок действия проверяется только для не-API запросов (isApi).
func (o OAuth2) IntrospectOAuth2Token(rawToken string, isApi bool) (*entity.UserInfo, error) {
	if rawToken == "" {
		return nil, errors.New("empty token")
	}
	if o.verifier == nil {
		return nil, errors.New("oidc provider is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idToken, err := o.verifyToken(ctx, rawToken)
	if err != nil {
		return nil, fmt.Errorf("token is not valid: %w", err)
	}

	if !isApi {
		if idToken.Expiry.IsZero() {
			return nil, errors.New("token has no expiry")
		}
		if time.Now().After(idToken.Expiry) {
			return nil, errors.New("token expired")
		}
	}

	user := entity.UserInfo{}
	if err := idToken.Claims(&user); err != nil {
		return nil, err
	}
	if user.Sub == "" {
		return nil, errors.New("token has no subject")
	}

	return &user, nil
}

// verifyToken проверяет токен основным приложением, а если он выдан другим
// приложением того же Authentik (например, Grafana) - verifier'ом этого приложения.
func (o OAuth2) verifyToken(ctx context.Context, rawToken string) (*oidc.IDToken, error) {
	idToken, err := o.verifier.Verify(ctx, rawToken)
	if err == nil {
		return idToken, nil
	}

	issuer := tokenIssuer(rawToken)
	if !o.isTrustedIssuer(issuer) {
		return nil, err
	}

	verifier, verr := o.verifierForIssuer(ctx, issuer)
	if verr != nil {
		return nil, err
	}
	return verifier.Verify(ctx, rawToken)
}

// verifierForIssuer возвращает (и кэширует) verifier для приложения Authentik.
func (o OAuth2) verifierForIssuer(ctx context.Context, issuer string) (*oidc.IDTokenVerifier, error) {
	if o.verifiers != nil {
		o.verifiers.mu.RLock()
		verifier, ok := o.verifiers.m[issuer]
		o.verifiers.mu.RUnlock()
		if ok {
			return verifier, nil
		}
	}

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}
	verifier := provider.Verifier(&oidc.Config{
		ClientID:          o.oAuth2Config.ClientID,
		SkipClientIDCheck: true,
		SkipExpiryCheck:   true,
	})

	if o.verifiers != nil {
		o.verifiers.mu.Lock()
		o.verifiers.m[issuer] = verifier
		o.verifiers.mu.Unlock()
	}
	return verifier, nil
}

// isTrustedIssuer разрешает issuer'ы основного приложения и других приложений
// того же экземпляра Authentik (тот же scheme+host, путь /application/o/<slug>/).
func (o OAuth2) isTrustedIssuer(issuer string) bool {
	if issuer == "" || o.mainIssuer == "" {
		return false
	}
	if issuer == o.mainIssuer {
		return true
	}

	mainURL, err := url.Parse(o.mainIssuer)
	if err != nil {
		return false
	}
	u, err := url.Parse(issuer)
	if err != nil {
		return false
	}
	if u.Scheme != mainURL.Scheme || !strings.EqualFold(u.Host, mainURL.Host) {
		return false
	}
	if u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return false
	}
	if !strings.HasPrefix(u.Path, authentikAppsPath) {
		return false
	}
	slug := strings.TrimSuffix(strings.TrimPrefix(u.Path, authentikAppsPath), "/")
	if slug == "" || strings.Contains(slug, "/") {
		return false
	}
	return issuerSlugRe.MatchString(slug)
}

// tokenIssuer достаёт iss из payload непроверенного JWT. Используется только
// для выбора verifier'а, сам токен всё равно проверяется по JWKS и issuer.
func tokenIssuer(rawToken string) string {
	parts := strings.Split(rawToken, ".")
	if len(parts) != 3 {
		return ""
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return ""
	}
	var claims struct {
		Issuer string `json:"iss"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ""
	}
	return claims.Issuer
}

// Revoke отзывает refresh-токен в Authentik (RFC 7009).
// Вызывается при выходе; отсутствие revocation endpoint не считается ошибкой.
func (o OAuth2) Revoke(refreshToken string) error {
	if refreshToken == "" || o.revocationURL == "" {
		return nil
	}

	form := url.Values{}
	form.Set("token", refreshToken)
	form.Set("token_type_hint", "refresh_token")
	form.Set("client_id", o.oAuth2Config.ClientID)
	form.Set("client_secret", o.oAuth2Config.ClientSecret)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.revocationURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("revocation endpoint returned %s", resp.Status)
	}

	return nil
}
