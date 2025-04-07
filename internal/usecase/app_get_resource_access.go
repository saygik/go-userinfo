package usecase

func (u *UseCase) GetAccessToResource(resource string, user string) (res int) {
	domain := getDomainFromUserName(user)
	res = -1
	userRoles, err := u.repo.GetUserRoles(user)
	if err != nil {
		return -1
	}
	if IsStringInArrayIdName("Администратор системы", userRoles) {
		return 1
	}
	if IsStringInArrayIdName("Администратор", userRoles) {
		return 1
	}
	if IsStringInArrayIdName("Технический специалист", userRoles) && domain == resource {
		return 1
	}

	accessRole, _ := u.repo.GetUserResourceAccess(resource, user)

	if accessRole != 1 && IsStringInArrayIdName("Группа компетенции БЖД", userRoles) {
		return 2
	}
	return accessRole
}
