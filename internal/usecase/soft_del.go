package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) DelUserSoftware(softwareForm entity.SoftwareForm) error {
	return u.repo.DelOneUserSoftware(softwareForm)
}
