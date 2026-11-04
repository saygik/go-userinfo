package usecase

import (
	"fmt"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

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

func (u *UseCase) DelUserSection(perms entity.Permissions, user string, section string) error {

	if !perms.IsSysAdmin {
		return u.Error("у вас нет прав для выполнения операции")
	}
	if user == "" {
		return u.Error("пользователь не определён")
	}
	if !isEmailValid(user) {
		return u.Error("неправильное имя пользователя")
	}
	err := u.repo.DelUserSection(user, section)
	if err != nil {
		return u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	return nil
}

func (u *UseCase) DelUserDomainRole(perms entity.Permissions, user string, domain string) error {

	if !perms.IsSysAdmin {
		return u.Error("у вас нет прав для выполнения операции")
	}
	if user == "" {
		return u.Error("пользователь не определён")
	}
	if !isEmailValid(user) {
		return u.Error("неправильное имя пользователя")
	}
	err := u.repo.DelUserDomainRole(user, domain)
	if err != nil {
		return u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	return nil
}

// DeleteAllUserRights удаляет все роли, доступы к доменам и дополнительные
// разделы пользователя. Текущая учётная запись защищена от саморазблокировки.
func (u *UseCase) DeleteAllUserRights(perms entity.Permissions, user string) error {

	if !perms.IsSysAdmin {
		return u.Error("у вас нет прав для выполнения операции")
	}
	if user == "" {
		return u.Error("пользователь не определён")
	}
	if !isEmailValid(user) {
		return u.Error("неправильное имя пользователя")
	}
	if perms.User != "" && strings.EqualFold(strings.TrimSpace(perms.User), strings.TrimSpace(user)) {
		return u.Error("нельзя удалить все права у своей учётной записи")
	}
	err := u.repo.DeleteAllUserRights(user)
	if err != nil {
		return u.Error(fmt.Sprintf("ошибка базы данных: %s", err.Error()))
	}
	return nil
}
