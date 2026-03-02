package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetCurrentUser(perms entity.Permissions) (map[string]interface{}, error) {
	user, err := u.GetUser(perms.User)
	if err != nil {
		return nil, u.Error(err.Error())
	}
	userG, err := u.GetGlpiUserForTechnical(perms.User, perms.User)
	if err == nil {
		user["glpi"] = userG
	}
	user["perms"] = perms
	return user, nil

}
