package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetSchedule(id int) (schedule entity.Schedule, err error) {
	err = r.db.SelectOne(&schedule, "GetSchedule $1", id)
	return schedule, err
}
