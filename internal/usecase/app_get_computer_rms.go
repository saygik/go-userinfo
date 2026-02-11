package usecase

func (u *UseCase) GetComputerRMS(user string) ([]string, error) {
	return u.repo.GetComputerRMS(user)
}
