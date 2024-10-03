package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetMattermostUsers() ([]entity.MattermostUserWithSessions, error) {
	//u.matt.GetUsersWithSessions()
	return u.matt.GetUsers()
}
