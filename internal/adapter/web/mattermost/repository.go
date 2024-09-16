package mattermost

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

type Repository struct {
	url    string
	token  string
	client *model.Client4
}

func New(url string, token string) *Repository {
	client := model.NewAPIv4Client("https://matt.rw")
	//	client.Login("bot@example.com", "password")
	client.SetToken("jmhmhqq87jr87edud897a8am3e")

	return &Repository{
		url:    url,
		token:  token,
		client: client,
	}
}
