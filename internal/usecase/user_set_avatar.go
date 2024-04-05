package usecase

import (
	"fmt"
)

func (u *UseCase) SetUserAvatar(userTechName string, user string, avatar string) error {

	if userTechName == "" {
		return u.Error("сначала войдите в систему")
	}
	if user == "" {
		return u.Error("неправильное имя пользователя")
	}
	if !isEmailValid(user) {
		return u.Error("неправильное имя пользователя")
	}
	err := u.repo.SetUserAvatar(user, avatar)
	if err != nil {
		return u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	return nil
}
