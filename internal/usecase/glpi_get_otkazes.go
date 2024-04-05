package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetGLPIOtkazes(startdate string, enddate string) (otkazes []entity.GLPI_Otkaz, err error) {
	return u.glpi.GetOtkazes(startdate, enddate)
}
