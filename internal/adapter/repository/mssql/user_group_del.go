package mssql

func (r *Repository) DelUserGroup(user string, id int) error {
	res, err := r.db.Exec("DelfindUserUserRole $1,$2", user, id)
	return sqlRowAffectedErrorWrapper(res, err)
}
