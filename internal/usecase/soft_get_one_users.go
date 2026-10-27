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
		return nil, u.Error("невозможно получить пользователей системы")
	}

	softUsers := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		softUsers = append(softUsers, u.softUserResponse(user))
	}

	return softUsers, nil
}
