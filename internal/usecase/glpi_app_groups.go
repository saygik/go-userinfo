package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetAppGroups(user string) (otkazes []entity.IdName, err error) {
	return u.repo.GetAppGroups()
}
