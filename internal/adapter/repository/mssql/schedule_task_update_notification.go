package mssql

func (r *Repository) UpdateScheduleTaskNotification(id int) error {
	res, err := r.db.Exec("UpdateScheduleTaskNotification $1", id)
	return sqlRowAffectedErrorWrapper(res, err)
}
