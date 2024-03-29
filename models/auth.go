package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/saygik/go-userinfo/db"
	uuid "github.com/twinj/uuid"
)

// TokenDetails ...
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

// AccessDetails ...
type AccessDetails struct {
	AccessUUID string
	UserID     string
}

// Token ...
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type OAuthTokens struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// User in Redis
type UserInRedis struct {
	Login string `json:"login"`
}

// User in Redis OpenID
type UserInRedisOpenID struct {
	Login       string `json:"login"`
	AccessToken string `json:"access_token"`
}

// OpenID User
type OpenIDUser struct {
	Id        string `json:"sub"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
}

// AuthModel ...
type AuthModel struct{}

// CreateToken ...
func (m AuthModel) CreateToken(userLogin string) (*TokenDetails, error) {

	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Hour * 4).Unix()
	//  td.AtExpires = time.Now().Add(time.Minute * 60).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	//td.RtExpires = time.Now().Add(time.Minute * 5).Unix()
	td.RefreshUUID = uuid.NewV4().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userLogin
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userLogin
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// 899063632300
// CreateAuth ...
func (m AuthModel) CreateAuth(userLogin string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	userInRedis := UserInRedis{
		Login: userLogin}
	//userInRedis.ID=user.ID
	//userInRedis.AO=user.AO
	userInRedisJSON, _ := json.Marshal(userInRedis)
	errAccess := db.GetRedis().Set(ctx, td.AccessUUID, userInRedisJSON, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := db.GetRedis().Set(ctx, td.RefreshUUID, userInRedisJSON, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

// ExtractToken ...
func (m AuthModel) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// VerifyToken ...
func (m AuthModel) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := m.ExtractToken(r)
	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	t := []byte("qWfD4_C8Q5T5y1Zx1_uQG8V9RDJUcPwWzyXDRGXL")
	// 	return t, nil
	// })
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
func (m AuthModel) VerifyTokenByTokenString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		//return []byte("qWfD4_C8Q5T5y1Zx1_uQG8V9RDJUcPwWzyXDRGXL"), nil
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
func (m AuthModel) ExtractTokenMetadataByTokenString(tokenString string) (*AccessDetails, error) {
	token, err := m.VerifyTokenByTokenString(tokenString)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID := fmt.Sprintf("%.f", claims["user_id"])
		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}

// TokenValid ...
func (m AuthModel) TokenValid(r *http.Request) error {
	token, err := m.VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// ExtractTokenMetadata ...
func (m AuthModel) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := m.VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID := claims["user_id"].(string)
		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}

// FetchAuth ...
func (m AuthModel) FetchAuth(authD *AccessDetails) (UserInRedis, error) {
	userJSON, err := db.GetRedis().Get(ctx, authD.AccessUUID).Result()
	user := UserInRedis{}
	if err != nil {
		return user, err
	}
	json.Unmarshal([]byte(userJSON), &user)
	//	userID, _ := strconv.ParseInt(userid, 10, 64)
	return user, nil
}

// DeleteAuth ...
func (m AuthModel) DeleteAuth(givenUUID string) (int64, error) {
	deleted, err := db.GetRedis().Del(ctx, givenUUID).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// CreateAuth ...
func (m AuthModel) CreateAuthOpenID(oauthToken OAuthTokens) error {

	td, err := m.ExtractTokenMetadataByTokenString(oauthToken.IdToken)

	//	token, err := m.VerifyTokenByTokenString(stoken)
	if err != nil {
		return err
	}

	at := time.Unix(int64(oauthToken.ExpiresIn), 0) //converting Unix to UTC(to Time object)
	now := time.Now()
	userInRedis := UserInRedisOpenID{
		Login:       td.UserID,
		AccessToken: oauthToken.AccessToken}
	userInRedisJSON, _ := json.Marshal(userInRedis)
	errAccess := db.GetRedis().Set(ctx, td.AccessUUID, userInRedisJSON, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	//	userJSON, err := db.GetRedis().Get(ctx, td.AccessUUID).Result()
	//	fmt.Printf(userJSON)
	return nil
}
func (m AuthModel) GetUserInfo(token string) (user OpenIDUser, err error) {
	req, _ := http.NewRequest("GET", "https://adss.brnv.rw/sso/oauth/7dd4815e19af5fbea99a290b134b7e493569ea13/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return user, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	if response.StatusCode != 200 {
		return user, fmt.Errorf("Access token is invalid or expired.")
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return user, fmt.Errorf("failed reading user info from server: %s", err.Error())
	}
	//	var user OpenIDUser

	json.Unmarshal([]byte(contents), &user)

	return user, nil

}
