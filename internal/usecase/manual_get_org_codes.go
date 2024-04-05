package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetOrgCodes() (orgs []entity.OrgWithCodes, err error) {
	return u.repo.GetOrgCodes()
}
