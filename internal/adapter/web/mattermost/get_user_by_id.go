package mattermost

import (
	"context"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetUserById(id string) (*entity.MattermostUser, error) {

	user, _, err := r.client.GetUser(context.Background(), id, "")
	if err != nil {
		return nil, err
	}
	mUser := entity.MattermostUser{
		Id:          user.Id,
		AuthService: user.AuthService,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Nickname:    user.Nickname,
		Name:        user.Username,
	}
	_ = user
	return &mUser, err
}
