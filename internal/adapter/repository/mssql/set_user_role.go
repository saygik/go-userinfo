package mssql

func (r *Repository) SetUserRole(user string, id int) error {
	res, err := r.db.Exec("UpdatefindUserUserRoles $1,$2", user, id)
	return sqlRowAffectedErrorWrapper(res, err)
}
