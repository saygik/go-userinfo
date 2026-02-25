package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) AddUserGroup(perms entity.Permissions, user string, id int) error {

	if !perms.IsSysAdmin {
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

func (u *UseCase) AddUserRole(perms entity.Permissions, user string, id int) error {

	if !perms.IsSysAdmin {
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
