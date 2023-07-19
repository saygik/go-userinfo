package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/forms"
	"github.com/saygik/go-userinfo/models"
)

// AuthController ...
type AuthController struct{}

var authModel = new(models.AuthModel)

// UserFromToken ...
func (ctl AuthController) UserFromToken(c *gin.Context) {

	tokenAuth, err := authModel.ExtractTokenMetadata(c.Request)
	if err != nil {
		//Token either expired or not valid
		return
	}

	user, err := authModel.FetchAuth(tokenAuth)
	if err != nil {
		//Token does not exists in Redis (User logged out or expired)
		return
	}

	//To be called from GetUserID()
	c.Set("user", user)
}

// TokenValid ...
func (ctl AuthController) TokenValid(c *gin.Context) {

	tokenAuth, err := authModel.ExtractTokenMetadata(c.Request)
	if err != nil {
		//Token either expired or not valid
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Cначала войдите в систему"})
		return
	}

	user, err := authModel.FetchAuth(tokenAuth)
	if err != nil {
		//Token does not exists in Redis (User logged out or expired)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Cначала войдите в систему"})
		return
	}

	//To be called from GetUserID()
	c.Set("user", user)
}

// Refresh ...
func (ctl AuthController) Refresh(c *gin.Context) {
	var tokenForm forms.Token

	if c.ShouldBindJSON(&tokenForm) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Invalid form", "form": tokenForm})
		c.Abort()
		return
	}

	//verify the token
	token, err := jwt.Parse(tokenForm.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
			return
		}
		userID := fmt.Sprintf("%s", claims["user_id"])
		//if err != nil {
		//	c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
		//	return
		//}
		//Delete the previous Refresh Token

		//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		deleted, delErr := authModel.DeleteAuth(refreshUUID)
		if delErr != nil || deleted == 0 { //if any goes wrong
			//c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
			//return
		}
		//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

		//Create new pairs of refresh and access tokens
		ts, createErr := authModel.CreateToken(userID)
		if createErr != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid authorization, please login again"})
			return
		}
		////Get user from BD
		//user, err := userModel.One(userID)
		//if err != nil {
		//	c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
		//	return
		//}
		//save the tokens metadata to redis
		saveErr := authModel.CreateAuth(userID, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid authorization, please login again"})
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusOK, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
	}
}

func (ctrl UserController) LoginCallback(c *gin.Context) {
	code := c.Query("code")
	redirectUri := c.Query("redirect_uri")
	if code == "" {
		code = DefaultDomain
	}
	tokens, err := getToken(code, redirectUri)
	//	fmt.Printf("%s\n", tokens)
	if err != nil {
		//		fmt.Println(err.Error())
		//		c.Redirect(http.StatusTemporaryRedirect, "http://10.2.146.202:3000")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
		return
	}
	authModel.CreateAuthOpenID(tokens)
	tokenRes := map[string]string{
		"access_token": tokens.IdToken,
	}
	c.JSON(http.StatusOK, tokenRes)
	// c.Redirect(http.StatusTemporaryRedirect, "http://10.2.146.202:3000/auth?token="+tokens.IdToken)
}

func getToken(token, redirectUri string) (models.OAuthTokens, error) {

	// if state != oauthStateString {
	// 	return nil, fmt.Errorf("invalid oauth state")
	// }

	// token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	// if err != nil {
	// 	return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	// }
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", "4OcmRaXcBIsoTehRDcF5fYO3N")
	data.Set("client_secret", "qWfD4_C8Q5T5y1Zx1_uQG8V9RDJUcPwWzyXDRGXL")
	data.Set("redirect_uri", redirectUri)
	data.Set("code", token)

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, "https://adss.brnv.rw/sso/oauth/7dd4815e19af5fbea99a290b134b7e493569ea13/token", strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(r)
	var requestToken models.OAuthTokens
	if err != nil {
		return requestToken, fmt.Errorf("failed getting access token: %s", err.Error())
	}
	if response.StatusCode != 200 {
		return requestToken, fmt.Errorf("failed getting access token from provider")
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)

	if err != nil {
		return requestToken, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	json.Unmarshal([]byte(contents), &requestToken)

	// req, _ := http.NewRequest("GET", "https://adss.brnv.rw/sso/oauth/7dd4815e19af5fbea99a290b134b7e493569ea13/userinfo", nil)
	// req.Header.Set("Authorization", "Bearer "+requestToken.AccessToken)
	// response, err = client.Do(req)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	// }
	// defer response.Body.Close()
	// contents, err = io.ReadAll(response.Body)

	return requestToken, nil
}
