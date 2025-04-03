package usecase

import (
	"strconv"
	"time"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) AddOneSoftwareUser(id string, softwareForm entity.SoftUser, editor string) (map[string]interface{}, error) {
	if idd, err := strconv.ParseInt(id, 10, 64); err != nil {
		return nil, u.Error("неправильный id пользователя")
	} else {
		softwareForm.Id = idd
	}
	currentTime := time.Now().Format("2006-01-02T15:04:05.000Z")
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

func (u *UseCase) UpdateOneSoftwareUser(softwareForm entity.SoftUser, editor string) error {
	currentTime := time.Now().Format("2006-01-02T15:04:05.000Z")
	//.Format("2006-01-02 15:04:05.000000 07:00")
	err := u.repo.UpdateOneSoftwareUser(softwareForm, editor, currentTime)
	if err != nil {
		return u.Error("ошибка SQL изменения пользователя в системе")
	}
	return nil
}
