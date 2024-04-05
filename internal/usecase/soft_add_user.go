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
	err := u.repo.AddOneSoftwareUser(softwareForm)
	if err != nil {
		return nil, u.Error("ошибка SQL добавления пользователя в систему")
	}
	userProperties := u.GetUserADPropertysShort(softwareForm.Name)
	userProperties["name"] = softwareForm.Name
	userProperties["login"] = softwareForm.Login
	userProperties["comment"] = softwareForm.Comment
	userProperties["fio"] = softwareForm.Fio
	userProperties["external"] = softwareForm.External
	return userProperties, nil
}
