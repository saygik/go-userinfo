package usecase

func (u *UseCase) SendPostSimple(channel, message string) error {

	return u.matt.SendPostSimple(channel, message)
}
func (u *UseCase) ReplyPost(channel, postId, message string) error {
	return u.matt.ReplyPost(channel, postId, message)
}
