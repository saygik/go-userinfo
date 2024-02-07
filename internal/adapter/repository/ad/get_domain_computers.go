package ad

func (r *Repository) GetDomainComputers(domain string) ([]map[string]interface{}, error) {
	return r.ads[domain].GetAllComputers()

}
