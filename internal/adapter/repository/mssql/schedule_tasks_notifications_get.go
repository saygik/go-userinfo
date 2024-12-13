package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetScheduleTasksNotifications(date string) (tasks []entity.ScheduleTask, err error) {
	_, err = r.db.Select(&tasks, "GetScheduleNotifications $1", date)
	return tasks, err
}
