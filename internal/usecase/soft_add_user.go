package usecase

import (
	"strconv"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) AddOneSoftwareUser(id string, softwareForm entity.SoftUser) (map[string]interface{}, error) {
	if idd, err := strconv.ParseInt(id, 10, 64); err != nil {
		return nil, u.Error("неправильный id пользователя")
	} else {
		softwareForm.Id = idd
	}
	user, err := u.repo.AddOneSoftwareUser(softwareForm)
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

	return adUser, nil
}
