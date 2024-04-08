package ad

func (r *Repository) GetGroupUsers(domain string, group string) ([]map[string]interface{}, error) {
	return r.ads[domain].GetGroupUsers(group)

}
