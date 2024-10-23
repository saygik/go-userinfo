package mattermost

import (
	"github.com/mattermost/mattermost/server/public/model"
)

type User *model.User
type Session *model.Session

type Repository struct {
	url    string
	token  string
	client *model.Client4
}

func New(url string, token string) *Repository {
	client := model.NewAPIv4Client("https://matt.rw")
	//	client.Login("bot@example.com", "password")
	client.SetToken(token)

	return &Repository{
		url:    url,
		token:  token,
		client: client,
	}
}

//czhom66383rc9ecdx669nnaaka
