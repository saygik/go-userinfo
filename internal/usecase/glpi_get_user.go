package usecase

import (
	"fmt"
	"strings"

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
	groups, err := u.repo.GetUserGlpiTrackingGroups(user)
	if err != nil {
		glpiUser.TrackingGroups = []entity.IdName{}
		return glpiUser, nil
	}

	str := strings.Replace(strings.Replace(strings.Replace(strings.Trim(fmt.Sprint(groups), "[]"), " ", ",", -1), "{", "", -1), "}", "", -1)
	if len(str) < 1 {
		glpiUser.TrackingGroups = []entity.IdName{}
		return glpiUser, nil
	}
	tg, err := u.glpi.GetUserTrackingGroups(str)
	if err != nil {
		glpiUser.TrackingGroups = []entity.IdName{}
		return glpiUser, nil
	}
	glpiUser.TrackingGroups = tg
	return glpiUser, nil
}
