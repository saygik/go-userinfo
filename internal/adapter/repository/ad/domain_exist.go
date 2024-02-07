package ad

func (r *Repository) IsDomainExist(domain string) bool {
	for _, oneADConfig := range r.ads {
		if oneADConfig.Domain == domain {
			return true
		}
	}
	return false
}
