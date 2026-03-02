package usecase

// GetADComputersOSFamily возвращает количество компьютеров домена, сгруппированное по семейству ОС (OperatingSystemFamily).
func (u *UseCase) ComputerLocalAdminsAudit(computer string, localAdmins []string, isDomain bool) error {
	return u.repo.ComputerLocalAdminsAudit(computer, localAdmins, isDomain)
}
