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
		adUser["login"] = user.Login
		adUser["comment"] = user.Comment
		adUser["fio"] = user.Fio
		adUser["external"] = user.External
		softUsers = append(softUsers, adUser)
	}

	return softUsers, nil

}
