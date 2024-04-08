package mssql

/********** Avatars*/

func (r *Repository) GetUserAvatar(userID string) (avatar string, err error) {
	err = r.db.QueryRow("GetUserAvatar $1", userID).Scan(&avatar)
	return avatar, err
}
