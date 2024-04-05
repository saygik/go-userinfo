package usecase

func (u *UseCase) IsSysAdmin(user string) bool {

	userRoles, err := u.repo.GetUserRoles(user)
	if err != nil {
		return false
	}
	if IsStringInArrayIdName("Администратор системы", userRoles) {
		return true
	}

	return false
}
