package mattermost

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetUsers() ([]entity.MattermostUser, error) {
	users := []entity.MattermostUser{}
	url_api := r.url
	urlUsers := url_api + "/users/search"

	var jsonStr = []byte(`{ "term": "*", "allow_inactive": false, "limit": 1000 }`)
	req, err := http.NewRequest("POST", urlUsers, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)

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
