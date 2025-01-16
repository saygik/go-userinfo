package usecase

func (u *UseCase) MattermostIntegrationAllowed(ip string) bool {
	hosts := u.matt.IntegrationAllowedHosts()
	return IsStringInArray(ip, hosts)
}
