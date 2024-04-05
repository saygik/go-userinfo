package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetGLPIProblems(startdate string, enddate string) ([]entity.GLPI_Problem, error) {
	return u.glpi.GetProblems(startdate, enddate)
}
