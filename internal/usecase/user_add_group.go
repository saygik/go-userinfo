package usecase

import (
	"fmt"
)

func (u *UseCase) AddUserGroup(techUser string, user string, id int) error {

	if isSysAdmin := u.IsSysAdmin(techUser); !isSysAdmin {
		return u.Error("у вас нет прав для выполнения операции")
	}
	if user == "" {
		return u.Error("пользователь не определён")
	}
	if !isEmailValid(user) {
		return u.Error("неправильное имя пользователя")
	}
	err := u.repo.AddUserGroup(user, id)
	if err != nil {
		return u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	return nil
}

func (u *UseCase) AddUserRole(techUser string, user string, id int) error {

	if isSysAdmin := u.IsSysAdmin(techUser); !isSysAdmin {
		return u.Error("у вас нет прав для выполнения операции")
	}
	if user == "" {
		return u.Error("пользователь не определён")
	}
	if !isEmailValid(user) {
		return u.Error("неправильное имя пользователя")
	}
	err := u.repo.AddUserRole(user, id)
	if err != nil {
		return u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	return nil
}
