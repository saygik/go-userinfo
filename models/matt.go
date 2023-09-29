package models

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// Sharepoint Model ...
type MattermostModel struct{}

type MattermostUser struct {
	Id          string `json:"id"`
	Name        string `json:"username"`
	AuthService string `json:"auth_service,omitempty"`
	AD          string `json:"ad,omitempty"`
	Email       string `json:"email"`
	Nickname    string `json:"nickname,omitempty"`
	Roles       string `json:"roles,omitempty"`
	IsBot       bool   `json:"is_bot,,omitempty"`
}

// GLPI User find by Mail ...
func (m MattermostModel) GetAll() (users []MattermostUser, err error) {
	url_api := os.Getenv("MATT_API")
	urlUsers := url_api + "/users/search"

	var jsonStr = []byte(`{ "term": "*", "allow_inactive": false, "limit": 1000 }`)
	req, err := http.NewRequest("POST", urlUsers, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("MATT_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return users, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &users)
	if err != nil {
		return users, err
	}

	return users, err
}
