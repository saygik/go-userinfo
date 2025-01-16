package ad

func (r *Repository) GetDomainAdminsGLPI(domain string) int {
	return r.adconfigs[domain].AdminGLPIGroup
}
