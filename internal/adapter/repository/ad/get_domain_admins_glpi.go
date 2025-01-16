package ad

func (r *Repository) GetDomainAdminsGLPI(domain string) int {
	if val, ok := r.adconfigs[domain]; ok {
		return val.AdminGLPIGroup
	} else {
		return 0
	}
}
