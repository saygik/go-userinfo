package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetGlpiUser(user string) (entity.GLPIUser, error) {
	glpiUser, err := u.glpi.GetUserByName(user)
	if err != nil {
		return entity.GLPIUser{}, u.Error("пользователь не найден в системе GLPI")
	}
	glpiUserProfiles, err := u.glpi.GetUserProfiles(glpiUser.Id)
	if err == nil {
		glpiUser.Profiles = glpiUserProfiles
	}
	glpiUserGroups, err := u.glpi.GetUserGroups(glpiUser.Id)
	if err == nil {
		glpiUser.Groups = glpiUserGroups
	}
	return glpiUser, nil
}
