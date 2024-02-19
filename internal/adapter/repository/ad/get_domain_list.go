package ad

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) DomainList() []entity.DomainList {
	var domainsNames []entity.DomainList
	for _, oneADConfig := range r.ads {
		domain := entity.DomainList{Name: oneADConfig.Domain, Title: oneADConfig.Title}
		domainsNames = append(domainsNames, domain)
	}
	return domainsNames
}
