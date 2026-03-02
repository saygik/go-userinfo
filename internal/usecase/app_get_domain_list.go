package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetDomainList(perms entity.Permissions) []entity.DomainList {

	domains := u.ad.DomainList()
	res := []entity.DomainList{}
	for _, oneDomain := range domains {
		accessLevel := u.GetAccessLevelForDomain(&perms, oneDomain.Name)
		if accessLevel != "none" {
			res = append(res, oneDomain)
		}

	}
	return res
}

func (u *UseCase) GetAccessLevelForDomain(perms *entity.Permissions, domainName string) string {
	// Приоритет: admin > tech > user
	if perms.AdminDomains[domainName] {
		return "admin"
	}
	if perms.TechDomains[domainName] {
		return "tech"
	}
	if perms.UserDomains[domainName] {
		return "user"
	}
	return "none"
}

func (u *UseCase) DomainList() []entity.DomainList {
	return u.ad.DomainList()
}
