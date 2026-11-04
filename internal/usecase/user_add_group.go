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
	if err != nil {
		return result, err
	}
	if result != nil {
		u.auditPermissionChange(perms, user, AuditRoleAdd, result.Name, "")
	}
	return result, nil
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
	if err != nil {
		return result, err
	}
	if result != nil {
		u.auditPermissionChange(perms, user, AuditSectionAdd, result.Name, "")
	}
	return result, nil
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
	if err != nil {
		return result, operation, err
	}
	action := AuditDomainAdd
	if operation == "UPDATED" {
		action = AuditDomainUpdate
	}
	if result != nil {
		u.auditPermissionChange(perms, user, action, result.Domain, result.AccessLevel)
	}
	return result, operation, nil
}
