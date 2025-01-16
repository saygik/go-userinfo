package mssql

func (r *Repository) UpdateScheduleTaskDisableMattermost(id int) error {
	res, err := r.db.Exec("UpdateScheduleTaskDisableMattermost $1", id)
	return sqlRowAffectedErrorWrapper(res, err)
}
