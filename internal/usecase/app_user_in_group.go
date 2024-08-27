package usecase

import (
	"strconv"
)

func (u *UseCase) UserInAppGroup(userID string, groupId string) error {
	groups, err := u.GetCurrentUserGroups(userID)
	if err != nil {
		return u.Error("Сприсок групп доступа приложения для данного пользователя пуст")
	}
	id, e := strconv.Atoi(groupId)
	if e != nil {
		return u.Error("ID группы доступа приложения не целое число")
	}
	for _, group := range groups {
		if group.Id == id {
			return nil
		}
	}

	return u.Error("пользователь отсутствует в группе приложения")
}
