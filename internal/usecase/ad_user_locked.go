package usecase

func (u *UseCase) ADUserLocked(user string) bool {
	domain := getDomainFromUserName(user)
	if !u.ad.IsDomainExist(domain) {
		return false
	}
	userinfo, err := u.ad.GetUser(domain, user)
	if err != nil {
		return false
	}
	uac := userinfo["userAccountControl"].(string)
	return (uac == "514" || uac == "66050")

}
