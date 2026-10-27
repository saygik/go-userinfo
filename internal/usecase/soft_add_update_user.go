package usecase

import (
	"strconv"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) AddOneSoftwareUser(id string, softwareForm entity.SoftUser, editor string) (map[string]interface{}, error) {
	softID, err := strconv.Atoi(id)
	if err != nil {
		return nil, u.Error("неправильный id системы")
	}
	softwareForm.IdSoft = softID
	softwareForm.Permissions = normalizeSoftUserPermissions(softwareForm.Permissions, true)

	currentTime, err := CurrentTimeFormattedRFC3339()
	if err != nil {
		return nil, u.Error("ошибка получения текущей даты и времени в формате RFC3339")
	}

	user, err := u.repo.AddOneSoftwareUser(softwareForm, editor, currentTime)
	if err != nil {
		return nil, u.Error("ошибка SQL добавления пользователя в систему")
	}

	return u.softUserResponse(user), nil
}

func (u *UseCase) UpdateOneSoftwareUser(softwareForm entity.SoftUser, editor string) (map[string]interface{}, error) {
	if softwareForm.Permissions != nil {
		softwareForm.Permissions = normalizeSoftUserPermissions(softwareForm.Permissions, false)
	}

	formatted, err := CurrentTimeFormattedRFC3339()
	if err != nil {
		return nil, u.Error("ошибка получения текущей даты и времени в формате RFC3339")
	}

	user, err := u.repo.UpdateOneSoftwareUser(softwareForm, editor, formatted)
	if err != nil {
		return nil, u.Error("ошибка SQL изменения пользователя в системе")
	}

	return u.softUserResponse(user), nil
}
