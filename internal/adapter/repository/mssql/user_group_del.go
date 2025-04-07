package mssql

func (r *Repository) DelUserGroup(user string, id int) error {
	res, err := r.db.Exec("DelfindUserUserGroup $1,$2", user, id)
	return sqlRowAffectedErrorWrapper(res, err)
}

func (r *Repository) DelUserRole(user string, id int) error {
	res, err := r.db.Exec("DelfindUserUserRole $1,$2", user, id)
	return sqlRowAffectedErrorWrapper(res, err)
}
