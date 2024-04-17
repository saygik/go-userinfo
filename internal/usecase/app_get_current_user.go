package usecase

import (
	"fmt"
)

func (u *UseCase) GetCurrentUser(userID string, techUser string) (map[string]interface{}, error) {
	user, err := u.GetUser(userID, techUser)
	if err != nil {
		return nil, u.Error(fmt.Sprintf("ошибка получения данных пользователя: %s", err.Error()))
	}
	userG, err := u.GetGlpiUserForTechnical(userID, techUser)
	if err == nil {
		user["glpi"] = userG
	}
	return user, nil

}
