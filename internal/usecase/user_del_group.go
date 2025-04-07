package usecase

import (
	"fmt"
)

func (u *UseCase) DelUserGroup(techUser string, user string, id int) error {

	if isSysAdmin := u.IsSysAdmin(techUser); !isSysAdmin {
		return u.Error("у вас нет прав для выполнения операции")
	}
	if user == "" {
		return u.Error("пользователь не определён")
	}
	if !isEmailValid(user) {
		return u.Error("неправильное имя пользователя")
	}
	err := u.repo.DelUserGroup(user, id)
	if err != nil {
		return u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	return nil
}

func (u *UseCase) DelUserRole(techUser string, user string, id int) error {

	if isSysAdmin := u.IsSysAdmin(techUser); !isSysAdmin {
		return u.Error("у вас нет прав для выполнения операции")
	}
	if user == "" {
		return u.Error("пользователь не определён")
	}
	if !isEmailValid(user) {
		return u.Error("неправильное имя пользователя")
	}
	err := u.repo.DelUserRole(user, id)
	if err != nil {
		return u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	return nil
}
