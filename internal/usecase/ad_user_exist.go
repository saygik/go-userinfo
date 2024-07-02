package usecase

func (u *UseCase) UserExist(userID string) error {

	domain := getDomainFromUserName(userID)
	if !u.ad.IsDomainExist(domain) {
		return u.Error("такой домен отсутствует в системе:")
	}
	_, err := u.redis.GetKeyFieldValue("allusers", userID)

	if u.ADUserLocked(userID) {
		return u.Error("пользователь существует, но заблокирован в соответствующем домене БЖД")
	}
	if err != nil {

		return u.Error("нет такого пользователя в зарегистрированных доменах БЖД")
	}

	return nil
}
