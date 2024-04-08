package mssql

func (r *Repository) DelScheduleTask(id string) error {
	res, err := r.db.Exec("DelScheduleTask $1", id)
	return sqlRowAffectedErrorWrapper(res, err)
}
