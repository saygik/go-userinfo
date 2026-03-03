package usecase

import "github.com/saygik/go-userinfo/internal/entity"

// GetADComputersOSFamily возвращает количество компьютеров домена, сгруппированное по семейству ОС (OperatingSystemFamily).
func (u *UseCase) ComputerLocalAdminsAudit(computer string, localAdmins []string, isDomain bool) error {
	return u.repo.ComputerLocalAdminsAudit(computer, localAdmins, isDomain)
}
func (u *UseCase) ComputerLocalAdminsGet(isDomain bool) (results []entity.LocalAdmins, err error) {
	return u.repo.ComputerLocalAdminsGet(isDomain)
}
