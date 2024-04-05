package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetMattermostUsers() ([]entity.MattermostUser, error) {
	return u.matt.GetUsers()
}
