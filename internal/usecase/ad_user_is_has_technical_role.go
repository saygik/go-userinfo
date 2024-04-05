package usecase

func (u *UseCase) HasTechnicalRole(user string) bool {
	userRoles, err := u.repo.GetUserRoles(user)
	if err != nil {
		return false
	}
	if IsStringInArrayIdName("Администратор системы", userRoles) {
		return true
	}
	if IsStringInArrayIdName("Администратор", userRoles) {
		return true
	}
	if IsStringInArrayIdName("Технический специалист", userRoles) {
		return true
	}
	return false
}
