package mattermost

func (r *Repository) IntegrationAllowedHosts() []string {
	return r.integrations.AllowedHosts
}
