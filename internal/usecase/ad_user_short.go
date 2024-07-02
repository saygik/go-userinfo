package usecase

func (u *UseCase) GetUserShort(userID string) (map[string]interface{}, error) {
	var user map[string]interface{}
	domain := getDomainFromUserName(userID)
	if !u.ad.IsDomainExist(domain) {
		return nil, u.Error("такой домен отсутствует в системе")
	}
	userJSON, err := u.redis.GetKeyFieldValue("allusers", userID)
	if err != nil {
		return nil, u.Error("нет такого пользователя в системе")
	}
	unmarshalString(userJSON, &user)
	user["domain"] = domain
	return user, nil
}
