package usecase

func (u *UseCase) GetMattermostUserById(id string) (string, string, error) {
	//u.matt.GetUsersWithSessions()
	mattUser, err := u.matt.GetUserById(id)
	if err != nil {
		return "", "", err
	}
	if len(mattUser.Nickname) > 0 {
		return mattUser.Nickname, mattUser.Name, nil
	} else {
		return mattUser.LastName + " " + mattUser.FirstName, mattUser.Name, nil
	}

}
