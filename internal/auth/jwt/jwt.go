package jwt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/saygik/go-userinfo/internal/entity"

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

// User
type UserAuth struct {
	Login string `json:"login"`
}

type Redis interface {
	AddKeyValue(string, interface{}, time.Duration) error
	GetKeyValue(string) (string, error)
	Delete(string) error
}

type JwtCfg struct {
	AccessSecret  string `binding:"required"`
	RefreshSecret string `binding:"required"`
	AtExpires     int    `binding:"required"`
	RtExpires     int    `binding:"required"`
}
type Auth struct {
	rc  Redis
	cfg JwtCfg
}

func New(redis Redis, cfg JwtCfg) *Auth {
	return &Auth{
		rc:  redis,
		cfg: cfg,
	}
}

func (m Auth) Login(login string) (token entity.Token, err error) {

	//Generate the JWT auth token

	tokenDetails, _ := m.CreateToken(login)

	saveErr := m.CreateAuth(login, tokenDetails)
	if saveErr == nil {
		token.AccessToken = tokenDetails.AccessToken
		token.RefreshToken = tokenDetails.RefreshToken
	}

	return token, nil
}

func (m Auth) CreateToken(userLogin string) (*TokenDetails, error) {

	td := &TokenDetails{}
	ate, _ := time.ParseDuration(fmt.Sprintf("%dh", m.cfg.AtExpires))
	rte, _ := time.ParseDuration(fmt.Sprintf("%dh", m.cfg.RtExpires))
	td.AtExpires = time.Now().Add(ate).Unix()
	//  td.AtExpires = time.Now().Add(time.Minute * 60).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(rte).Unix()
	//	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
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
	td.AccessToken, err = at.SignedString([]byte(m.cfg.AccessSecret))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userLogin
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(m.cfg.RefreshSecret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// CreateAuth ...
func (m Auth) CreateAuth(userLogin string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	user := UserAuth{
		Login: userLogin}
	userJSON, _ := json.Marshal(user)
	//errAccess := db.GetRedis().Set(ctx, td.AccessUUID, userJSON, at.Sub(now)).Err()
	errAccess := m.rc.AddKeyValue(td.AccessUUID, userJSON, at.Sub(now))

	if errAccess != nil {
		return errAccess
	}
	//	errRefresh := db.GetRedis().Set(ctx, td.RefreshUUID, userJSON, rt.Sub(now)).Err()
	errRefresh := m.rc.AddKeyValue(td.RefreshUUID, userJSON, rt.Sub(now))
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

// DeleteAuth ...
func (m Auth) DeleteAuth(givenUUID string) error {
	return m.rc.Delete(givenUUID)
}

// ExtractToken ...
func (m Auth) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// ExtractTokenMetadata ...
func (m Auth) ExtractTokenMetadata(r *http.Request) (*entity.AccessDetails, error) {
	tokenString := m.ExtractToken(r)
	token, err := m.VerifyToken(tokenString)
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
		return &entity.AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}

// VerifyToken ...
func (m Auth) VerifyToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.cfg.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// FetchAuth ...
func (m Auth) FetchAuth(authD *entity.AccessDetails) (string, error) {
	userJSON, err := m.rc.GetKeyValue(authD.AccessUUID)
	user := UserAuth{}
	if err != nil {
		return "", err
	}
	json.Unmarshal([]byte(userJSON), &user)
	//	userID, _ := strconv.ParseInt(userid, 10, 64)
	return user.Login, nil
}

func (m Auth) RefreshToken(tokenString string) (map[string]string, error) {
	var tokens map[string]string
	//verify the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.cfg.RefreshSecret), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		return tokens, fmt.Errorf("invalid authorization, please login again")
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return tokens, fmt.Errorf("invalid authorization, please login again")
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			return tokens, fmt.Errorf("invalid authorization, please login again")
		}
		userID := fmt.Sprintf("%s", claims["user_id"])
		//if err != nil {
		//	c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
		//	return
		//}
		//Delete the previous Refresh Token

		//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		delErr := m.DeleteAuth(refreshUUID)
		if delErr != nil { //if any goes wrong
			return tokens, fmt.Errorf("invalid authorization, please login again")

			//c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
			//return
		}
		//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

		//Create new pairs of refresh and access tokens
		ts, createErr := m.CreateToken(userID)
		if createErr != nil {
			return tokens, fmt.Errorf("invalid authorization, please login again")

		}

		saveErr := m.CreateAuth(userID, ts)
		if saveErr != nil {
			return tokens, fmt.Errorf("invalid authorization, please login again")

		}
		tokens = map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return tokens, nil
	} else {
		return tokens, fmt.Errorf("invalid authorization, please login again")
	}
}
