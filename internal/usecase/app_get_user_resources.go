package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetCurrentUserResources(user string) ([]entity.AppResource, error) {
	return u.repo.GetCurrentUserResources(user)
}
