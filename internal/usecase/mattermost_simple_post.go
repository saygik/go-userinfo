package usecase

func (u *UseCase) SendPostSimple(channel, message string) error {

	return u.matt.SendPostSimple(channel, message)
}
