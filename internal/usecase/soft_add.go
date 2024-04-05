package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) AddUserSoftware(softwareForm entity.SoftwareForm) error {
	return u.repo.AddOneUserSoftware(softwareForm)
}
