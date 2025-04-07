package mssql

func (r *Repository) AddUserGroup(user string, id int) error {
	res, err := r.db.Exec("AddfindUserUserGroup $1,$2", user, id)
	return sqlRowAffectedErrorWrapper(res, err)
}
func (r *Repository) AddUserRole(user string, id int) error {
	res, err := r.db.Exec("AddfindUserUserRole $1,$2", user, id)
	return sqlRowAffectedErrorWrapper(res, err)
}
