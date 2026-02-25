package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetCurrentUser(userID string, perms entity.Permissions) (map[string]interface{}, error) {
	user, err := u.GetUser(userID, perms)
	if err != nil {
		return nil, u.Error(err.Error())
	}
	userG, err := u.GetGlpiUserForTechnical(userID, perms.User)
	if err == nil {
		user["glpi"] = userG
	}
	return user, nil

}
