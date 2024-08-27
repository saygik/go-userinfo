package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetCurrentUserGroups(user string) ([]entity.IdName, error) {
	return u.repo.GetUserGroups(user)
}
