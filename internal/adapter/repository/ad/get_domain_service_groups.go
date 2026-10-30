package ad

// GetDomainServiceGroups возвращает distinguishedName сервисных групп домена.
func (r *Repository) GetDomainServiceGroups(domain string) []string {
	if val, ok := r.adconfigs[domain]; ok {
		groups := make([]string, 0, len(val.ServiceGroups))
		for _, group := range val.ServiceGroups {
			if group != "" {
				groups = append(groups, group)
			}
		}
		return groups
	}
	return []string{}
}
