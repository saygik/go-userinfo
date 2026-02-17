package ad

func (r *Repository) GetDomainAdminsGLPI(domain string) int {
	if val, ok := r.adconfigs[domain]; ok {
		return val.AdminGLPIGroup
	} else {
		return 0
	}
}

//      "mattermostLogChannel": "gqnnu5asjfdhzkysbrqxor1hde",
func (r *Repository) GetDomainMattermostLogChannel(domain string) string {
	if val, ok := r.adconfigs[domain]; ok {
		return val.MattermostLogChannel
	} else {
		return ""
	}
}
