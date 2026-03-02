package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetUserRoles(user string) (otkazes []entity.IdNameDescription, err error) {
	return u.repo.GetUserRoles(user)
}

func (u *UseCase) GetAppRoles() ([]entity.IdNameDescription, error) {
	return u.repo.GetAppRoles()
}

func (u *UseCase) GetAppSections() ([]entity.IdNameDescription, error) {
	return u.repo.GetAppSections()
}
func (u *UseCase) GetSections(user string) ([]entity.IdNameDescription, error) {
	return u.repo.GetSections(user)
}
func (u *UseCase) GetDomainAccess(user string) ([]entity.DomainAccess, error) {
	return u.repo.GetDomainAccess(user)
}
