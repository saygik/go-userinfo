package mattermost

import (
	"github.com/mattermost/mattermost/server/public/model"
)

type User *model.User
type Session *model.Session

type Integrations struct {
	AddCommentFromApi                  string   `json:"add-comment-from-api,omitempty"`
	DisableCalendarTaskNotificationApi string   `json:"disable-calendar-task-notification-api,omitempty"`
	AllowedHosts                       []string `json:"allowed-hosts,omitempty"`
}

type Repository struct {
	url          string
	token        string
	integrations Integrations
	client       *model.Client4
}

func New(url string, token string, apiIntegrations Integrations) *Repository {
	client := model.NewAPIv4Client("https://matt.rw")
	//	client.Login("bot@example.com", "password")
	client.SetToken(token)

	return &Repository{
		url:          url,
		token:        token,
		client:       client,
		integrations: apiIntegrations,
	}
}

//czhom66383rc9ecdx669nnaaka
