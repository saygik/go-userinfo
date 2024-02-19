package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetDomainList(user string) []entity.DomainList {
	domain := getDomainFromUserName(user)
	domains := u.ad.DomainList()
	res := []entity.DomainList{}
	for _, oneDomain := range domains {
		access := u.GetAccessToResource(oneDomain.Name, user)
		if access == -1 && domain == oneDomain.Name {
			access = 0
		}
		if access == -1 {
			continue
		}
		if access != -1 {
			res = append(res, oneDomain)
		}

	}
	return res
}
