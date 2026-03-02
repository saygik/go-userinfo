package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) AddUserRole(perms entity.Permissions, user string, id int) (*entity.IdNameDescription, error) {

	if !perms.IsSysAdmin {
		return nil, u.Error("у вас нет прав для выполнения операции")
	}
	if user == "" {
		return nil, u.Error("пользователь не определён")
	}
	if !isEmailValid(user) {
		return nil, u.Error("неправильное имя пользователя")
	}
	result, err := u.repo.AddUserRole(user, id)
	return result, err
}

func (u *UseCase) AddUserSection(perms entity.Permissions, user string, section string) (*entity.IdNameDescription, error) {

	if !perms.IsSysAdmin {
		return nil, u.Error("у вас нет прав для выполнения операции")
	}
	if user == "" {
		return nil, u.Error("пользователь не определён")
	}
	if !isEmailValid(user) {
		return nil, u.Error("неправильное имя пользователя")
	}
	result, err := u.repo.AddUserSection(user, section)
	return result, err
}

func (u *UseCase) AddUserDomainRole(perms entity.Permissions, user string, domain string, level string) (*entity.DomainAccess, string, error) {

	if !perms.IsSysAdmin {
		return nil, "", u.Error("у вас нет прав для выполнения операции")
	}
	if user == "" {
		return nil, "", u.Error("пользователь не определён")
	}
	if !isEmailValid(user) {
		return nil, "", u.Error("неправильное имя пользователя")
	}
	result, operation, err := u.repo.AddUserDomainRole(user, domain, level)
	return result, operation, err
}
