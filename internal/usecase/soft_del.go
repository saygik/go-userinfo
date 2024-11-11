package usecase

func (u *UseCase) DelUserSoftware(id string) error {
	return u.repo.DelOneUserSoftware(id)
}
