package ad

func (r *Repository) GetDomainRMSPort(domain string) int {
	if val, ok := r.adconfigs[domain]; ok {
		if val.RmsPort == 0 {
			return 25650
		}
		return val.RmsPort
	} else {
		return 25650
	}
}
