package mssql

func (r *Repository) SetUserAvatar(user string, avatar string) error {
	res, err := r.db.Exec("SetUserAvatar $1,$2", user, avatar)
	return sqlRowAffectedErrorWrapper(res, err)
}
