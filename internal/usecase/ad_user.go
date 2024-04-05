package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetUser(userID string, techUser string) (map[string]interface{}, error) {
	var user map[string]interface{}
	domain := getDomainFromUserName(userID)
	if !u.ad.IsDomainExist(domain) {
		return nil, u.Error("такой домен отсутствует в системе")
	}
	userJSON, err := u.redis.GetKeyFieldValue("allusers", userID)
	if err != nil {
		return nil, u.Error("нет таккого пользователя в системе")
	}

	unmarshalString(userJSON, &user)

	// role := entity.IdName{Id: 5, Name: "Пользователь"}
	// roles, err := u.repo.GetUserRoles(userID)
	// if err == nil && len(roles) > 0 {
	// 	role = roles[0]
	// }
	// user["role"] = role
	user["domain"] = domain
	// groups, err := u.repo.GetUserGroups(userID)
	// if err == nil && len(groups) > 0 {
	// 	user["groups"] = groups
	// } else {
	// 	user["groups"] = []entity.IdName{}
	// }

	avatar, err := u.repo.GetUserAvatar(userID)
	if err == nil {
		user["avatar"] = avatar
	}
	if isSysAdmin := u.IsSysAdmin(techUser); !(isSysAdmin || userID == techUser) {
		return user, nil
	}
	userRole := u.repo.GetUserRole(userID)
	user["app_role"] = userRole

	groups, err := u.repo.GetUserGroups(userID)
	if err == nil && len(groups) > 0 {
		user["app_groups"] = groups
	} else {
		user["app_groups"] = []entity.IdName{}
	}

	return user, nil
}
