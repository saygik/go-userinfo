package mssql

func (r *Repository) DelOneUserSoftware(id string) (err error) {
	res, err := r.db.Exec(`
		DELETE FROM dbo.UserSoftwares
		WHERE id = $1`, id)
	return sqlRowAffectedErrorWrapper(res, err)
}
