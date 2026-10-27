package mssql

func (r *Repository) SetOneUserSoftwarePermissionSended(id int64) (err error) {
	res, err := r.db.Exec(`
		UPDATE dbo.UserSoftwaresPermissions
		SET sended = 1
		WHERE id = $1`, id)
	return sqlRowAffectedErrorWrapper(res, err)
}
