package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetScheduleTasks(id string) (tasks []entity.ScheduleTask, err error) {
	_, err = r.db.Select(&tasks, "GetScheduleTasks $1", id)
	return tasks, err
}
