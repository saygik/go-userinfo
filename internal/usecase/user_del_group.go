package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) DelUserGroup(perms entity.Permissions, user string, id int) error {

	if !perms.IsSysAdmin {
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

func (u *UseCase) DelUserRole(perms entity.Permissions, user string, id int) error {

	if !perms.IsSysAdmin {
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
