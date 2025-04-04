package usecase

import (
	"strconv"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) AddOneSoftwareUser(id string, softwareForm entity.SoftUser, editor string) (map[string]interface{}, error) {
	if idd, err := strconv.ParseInt(id, 10, 64); err != nil {
		return nil, u.Error("неправильный id пользователя")
	} else {
		softwareForm.Id = idd
	}
	currentTime, err := CurrentTimeFormattedRFC3339()
	if err != nil {
		return nil, u.Error("ошибка получения текущей даты и времени в формате RFC3339")
	}
	user, err := u.repo.AddOneSoftwareUser(softwareForm, editor, currentTime)
	if err != nil {
		return nil, u.Error("ошибка SQL добавления пользователя в систему")
	}

	adUser := u.GetUserADPropertysShort(user.Name)
	adUser["id"] = user.Id
	adUser["login"] = user.Login
	adUser["comment"] = user.Comment
	adUser["fio"] = user.Fio
	adUser["external"] = user.External
	if len(user.Mail) > 0 {
		adUser["mail"] = user.Mail
	}
	adUser["sended"] = user.Sended
	adUser["enddate"] = user.EndDate
	adUser["editor"] = user.Editor
	adUser["last_changes"] = user.LastChanges

	return adUser, nil
}

func (u *UseCase) UpdateOneSoftwareUser(softwareForm entity.SoftUser, editor string) (entity.SoftUser, error) {
	// Load the Europe/Minsk time zone
	formatted, err := CurrentTimeFormattedRFC3339()
	if err != nil {
		return entity.SoftUser{}, u.Error("ошибка получения текущей даты и времени в формате RFC3339")
	}

	user, err := u.repo.UpdateOneSoftwareUser(softwareForm, editor, formatted)
	if err != nil {
		return entity.SoftUser{}, u.Error("ошибка SQL изменения пользователя в системе")
	}
	return user, nil
}
