package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (h *Handler) GetLogin(c *gin.Context) {

	loginChallenge := c.Query("login_challenge")

	if loginChallenge == "" {
		//открыть форму с ошибкой return "", errors.New("the login_challenge parameter is present but empty")
		c.HTML(http.StatusOK, "login.html", gin.H{
			"ErrorTitle":   "Ошибка процесса входа в систему!",
			"ErrorContent": "Неавторизованный доступ",
		})
		return
	}

	resp, err := h.hydra.GetOAuth2LoginRequest(loginChallenge)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.GetOAuth2LoginRequest``: %v\n", err)
	}
	if resp.Skip {
		err := h.uc.UserExist(resp.Subject)
		if err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"ErrorTitle":   "Ошибка входа пользователя",
				"ErrorContent": err.Error(),
			})
			return
		}
		redirectTo, err := h.hydra.AcceptOAuth2LoginRequest(loginChallenge, resp.Subject)
		if err != nil {
			//открыть форму с ошибкой return "", errors.New("the login_challenge parameter is present but empty")

			c.HTML(http.StatusOK, "login.html", gin.H{
				"ErrorTitle":   "Невозможно подтвердить процесс входа в систему",
				"ErrorContent": err.Error(),
			})
			return
		}
		//return c.Redirect(http.StatusFound, resp2.RedirectTo)
		c.Redirect(http.StatusMovedPermanently, redirectTo)
		return
	}
	// return c.Render(http.StatusOK, "login.html", map[string]interface{}{
	// 	"LoginChallenge": loginChallenge,
	// })
	c.HTML(http.StatusOK, "login.html", gin.H{
		"LoginChallenge": loginChallenge,
	})
}

func (h *Handler) PostLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Выход произведен"})
}
func (h *Handler) GetLogout(c *gin.Context) {
	logoutChallenge := c.Query("logout_challenge")

	if logoutChallenge == "" {
		//открыть форму с ошибкой return "", errors.New("the login_challenge parameter is present but empty")
		c.HTML(http.StatusOK, "login.html", gin.H{
			"ErrorTitle":   "Login Challenge Is Not Exist!",
			"ErrorContent": "Login Challenge Is Not Exist!",
		})
		return
	}

	redirectTo, err := h.hydra.AcceptOAuth2LogoutRequest(logoutChallenge)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OAuth2API.GetOAuth2LoginRequest``: %v\n", err)
	}
	c.Redirect(http.StatusMovedPermanently, redirectTo)
}

func (h *Handler) PostLogin(c *gin.Context) {
	formData := struct {
		LoginChallenge string `validate:"required"`
		Email          string `validate:"required"`
		Password       string `validate:"required"`
		RememberMe     string `validate:"required"`
	}{
		LoginChallenge: c.PostForm("login_challenge"),
		Email:          c.PostForm("email"),
		Password:       c.PostForm("password"),
		RememberMe:     c.PostForm("remember_me"),
	}
	var rememberMe = formData.RememberMe == "true"
	loginForm := entity.LoginForm{
		Email:    formData.Email,
		Password: formData.Password,
	}
	err := h.uc.UserExist(loginForm.Email)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"ErrorTitle":   "Ошибка входа пользователя",
			"ErrorContent": err.Error(),
		})
		return
	}
	_, _, err = h.uc.Authenticate(loginForm)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"ErrorTitle":   "Неправильный логин или пароль",
			"ErrorContent": err.Error(),
		})
		return
	}

	subject := fmt.Sprint(formData.Email)
	redirect, err := h.hydra.AcceptNewOAuth2LoginRequest(formData.LoginChallenge, subject, rememberMe)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"ErrorTitle":   "Невозможно подтвердить запрос OAuth2 на авторизацию или истекло время ожидания",
			"ErrorContent": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusFound, redirect)

}

func (h *Handler) GetConsent(c *gin.Context) {

	consentChallenge := c.Query("consent_challenge")

	if consentChallenge == "" {
		c.HTML(http.StatusOK, "consent.html", gin.H{
			"ErrorTitle":   "Ошибка подтверждения входа в систему",
			"ErrorContent": "неправильный или пустой процесс входа",
		})
		return
	}

	resp, err := h.hydra.GetOAuth2ConsentRequest(consentChallenge)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Cначала войдите в систему"})
		return
	}

	//* проверка присутствия в необходимых scopes
	scopes := resp.RequestedScope
	err = h.uc.UserInGropScopes(*resp.Subject, scopes, h.hydra.GetScopes())
	if err != nil {
		c.HTML(http.StatusOK, "consent.html", gin.H{
			"ErrorTitle":   "Ошибка подтверждения входа в систему",
			"ErrorContent": err.Error(),
		})
		return
	}

	user, err := h.uc.GetUser(*resp.Subject, *resp.Subject)
	if err != nil {
		c.HTML(http.StatusOK, "consent.html", gin.H{
			"ErrorTitle":   "Ошибка получения данных пользователя",
			"ErrorContent": err.Error(),
		})
	}

	if *resp.Skip {
		redirectTo, err := h.hydra.AcceptOAuth2ConsentRequest(resp, user)
		if err != nil {
			c.HTML(http.StatusOK, "consent.html", gin.H{
				"ErrorTitle":   "Ошибка подтверждения входа в систему",
				"ErrorContent": err.Error(),
			})
			return
		}
		//return c.Redirect(http.StatusFound, resp2.RedirectTo)
		c.Redirect(http.StatusMovedPermanently, redirectTo)
	}

	return

}

func (h *Handler) PostConsent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Вход произведен"})
}
