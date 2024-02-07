package ad

func (r *Repository) GetDomainUsers(domain string) ([]map[string]interface{}, error) {
	return r.ads[domain].GetAllUsers()

}
