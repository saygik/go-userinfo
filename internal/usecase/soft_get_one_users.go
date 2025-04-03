package usecase

import (
	"strconv"
)

func (u *UseCase) GetSoftwareUsers(ids string) ([]map[string]interface{}, error) {
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, u.Error("неправильный id системы")
	}

	users, err := u.repo.GetSoftwareUsers(id)

	if err != nil {
		return nil, u.Error("невозможно получить систему из GLPI")
	}

	softUsers := []map[string]interface{}{}
	for _, user := range users {
		adUser := u.GetUserADPropertysShort(user.Name)
		adUser["id"] = user.Id
		adUser["login"] = user.Login
		adUser["comment"] = user.Comment
		adUser["fio"] = user.Fio
		adUser["external"] = user.External
		adUser["editor"] = user.Editor
		adUser["last_changes"] = user.LastChanges
		if len(user.Mail) > 0 {
			adUser["mail"] = user.Mail
		} else {
			adUser["admail"] = adUser["mail"]
			adUser["mail"] = ""
		}
		adUser["sended"] = user.Sended
		if user.EndDate[0:4] == "1900" {
			adUser["enddate"] = ""
		} else {
			adUser["enddate"] = user.EndDate
		}

		softUsers = append(softUsers, adUser)
	}

	return softUsers, nil

}
