package glpiapi

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GLPIInitSessionUser(usertoken string) (token string, err error) {
	glpirest := r.url + "initSession/"
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}

	req, err := http.NewRequest("GET", glpirest, nil)
	if err != nil {
		return token, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "user_token "+usertoken)
	req.Header.Set("App-Token", r.token)
	resp, err := client.Do(req)
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()

	var glpiToken entity.GlpiApiToken
	if resp.StatusCode == http.StatusOK {
		err1 := json.NewDecoder(resp.Body).Decode(&glpiToken)
		fmt.Println(err1)
		if err1 != nil {
			return token, err1
		}
		token = glpiToken.Token
	} else {
		return token, errors.New("невозможно создать новую сессию GLPI")
	}
	return token, err
}
