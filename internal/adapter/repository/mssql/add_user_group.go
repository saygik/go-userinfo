package mssql

func (r *Repository) AddUserGroup(user string, id int) error {
	res, err := r.db.Exec("AddfindUserUserRole $1,$2", user, id)
	return sqlRowAffectedErrorWrapper(res, err)
}
