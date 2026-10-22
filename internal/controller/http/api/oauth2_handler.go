package api

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	refreshTokenCookie = "refresh_token"
	oauthStateCookie   = "oauth2_auth_state"

	// Cookie авторизации нужны только роутам /api/oauth-authentik.
	authCookiePath = "/api/oauth-authentik"

	// Время жизни refresh-токена (в Authentik по умолчанию 30 дней).
	refreshTokenMaxAge = 30 * 24 * 60 * 60
	// Время на прохождение входа: state + PKCE code_verifier.
	oauthStateMaxAge = 10 * 60
)

func randomURLSafe(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// authCookieSecure - Secure обязателен для https и ломает вход по http (локальная разработка).
func (h *Handler) authCookieSecure() bool {
	redirect, err := url.Parse(h.oAuth2Authentik.GetRedirectUrl())
	return err == nil && redirect.Scheme == "https"
}

func (h *Handler) setAuthCookie(c *gin.Context, name, value string, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     authCookiePath,
		MaxAge:   maxAge,
		Secure:   h.authCookieSecure(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) deleteAuthCookie(c *gin.Context, name string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     authCookiePath,
		MaxAge:   -1,
		Expires:  time.Unix(1, 0).UTC(),
		Secure:   h.authCookieSecure(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) clearAuthCookies(c *gin.Context) {
	h.deleteAuthCookie(c, refreshTokenCookie)
	h.deleteAuthCookie(c, oauthStateCookie)
}

// LoginOauth2Authentik готовит ссылку на вход. state и PKCE code_verifier
// запоминаются на сервере в HttpOnly cookie, поэтому state проверяется сервером.
func (h *Handler) LoginOauth2Authentik(c *gin.Context) {
	state, err := randomURLSafe(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось создать состояние авторизации"})
		return
	}
	codeVerifier, err := randomURLSafe(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось создать состояние авторизации"})
		return
	}

	h.setAuthCookie(c, oauthStateCookie, state+"."+codeVerifier, oauthStateMaxAge)

	loginURL := h.oAuth2Authentik.AuthCodeURL(state, h.oAuth2Authentik.CodeChallenge(codeVerifier))
	c.JSON(http.StatusOK, gin.H{"data": loginURL})
}

// ExchangeTokenAuthentik меняет код авторизации на токены.
func (h *Handler) ExchangeTokenAuthentik(c *gin.Context) {
	var form struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		form.Code = c.Query("code")
		form.State = c.Query("state")
	}

	if form.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Не получен код авторизации"})
		return
	}

	stateCookie, err := c.Cookie(oauthStateCookie)
	if err != nil || stateCookie == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Сессия авторизации не найдена, начните вход заново"})
		return
	}

	expectedState, codeVerifier, ok := strings.Cut(stateCookie, ".")
	if !ok || expectedState == "" || codeVerifier == "" {
		h.deleteAuthCookie(c, oauthStateCookie)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректная сессия авторизации, начните вход заново"})
		return
	}
	if form.State == "" || form.State != expectedState {
		h.deleteAuthCookie(c, oauthStateCookie)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Проверка состояния авторизации не пройдена"})
		return
	}
	h.deleteAuthCookie(c, oauthStateCookie)

	token, user, err := h.oAuth2Authentik.Exchange(form.Code, codeVerifier)
	if err != nil {
		h.log.Errorf("authentik: обмен кода авторизации: %v", err)
		// 400, а не 401: клиент не должен пытаться обновить токен по этому ответу.
		c.JSON(http.StatusBadRequest, gin.H{"message": "Не удалось выполнить вход через Authentik", "error": err.Error()})
		return
	}

	if token.RefreshToken == "" {
		h.log.Warn("authentik: провайдер не вернул refresh_token - проверьте scope offline_access")
	} else {
		h.setAuthCookie(c, refreshTokenCookie, token.RefreshToken, refreshTokenMaxAge)
	}

	// Токены, которые доступны только серверу, в ответ клиенту не отдаём.
	token.RefreshToken = ""
	token.IdToken = ""

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// RefreshOauth2Authentik обновляет access-токен по refresh-токену из cookie.
func (h *Handler) RefreshOauth2Authentik(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil || refreshToken == "" {
		h.clearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Требуется повторный вход"})
		return
	}

	newToken, err := h.oAuth2Authentik.ExchangeRefreshToAccessToken(refreshToken)
	if err != nil {
		h.log.Warnf("authentik: обновление токена: %v", err)
		// Битый или отозванный токен не должен оставаться в браузере.
		h.clearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Сессия истекла, требуется повторный вход"})
		return
	}

	value := newToken.RefreshToken
	if value == "" {
		value = refreshToken
	}
	h.setAuthCookie(c, refreshTokenCookie, value, refreshTokenMaxAge)

	newToken.RefreshToken = ""
	newToken.IdToken = ""

	c.JSON(http.StatusOK, newToken)
}

// LogoutOauth2Authentik отзывает refresh-токен, чистит cookie и возвращает
// ссылку на завершение сессии Authentik с возвратом в приложение.
func (h *Handler) LogoutOauth2Authentik(c *gin.Context) {
	var form struct {
		Url string `json:"url,omitempty"`
	}
	_ = c.ShouldBindJSON(&form)

	if refreshToken, err := c.Cookie(refreshTokenCookie); err == nil && refreshToken != "" {
		if err := h.oAuth2Authentik.Revoke(refreshToken); err != nil {
			h.log.Warnf("authentik: отзыв refresh-токена: %v", err)
		}
	}
	h.clearAuthCookies(c)

	logoutURL := h.oAuth2Authentik.LogOutURL()
	if parsed, err := url.Parse(logoutURL); err == nil {
		query := parsed.Query()
		if form.Url != "" {
			query.Set("next", form.Url)
		}
		parsed.RawQuery = query.Encode()
		logoutURL = parsed.String()
	}

	c.JSON(http.StatusOK, gin.H{"data": logoutURL})
}
