package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) GetComputerRMS(user string) ([]entity.ComputerRMS, error) {
	return u.repo.GetComputerRMS(user)
}
