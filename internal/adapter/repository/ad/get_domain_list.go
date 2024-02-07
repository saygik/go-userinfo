package ad

func (r *Repository) DomainList() []string {
	var domainsNames []string
	for _, oneADConfig := range r.ads {
		domainsNames = append(domainsNames, oneADConfig.Domain)
	}
	return domainsNames
}
