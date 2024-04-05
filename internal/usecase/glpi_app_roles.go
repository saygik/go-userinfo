package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetAppRoles(user string) (otkazes []entity.IdName, err error) {
	return u.repo.GetAppRoles()
}
