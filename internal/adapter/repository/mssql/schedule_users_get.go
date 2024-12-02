package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetScheduleUsers(id string, tip int) (users []entity.IdName, err error) {
	_, err = r.db.Select(&users, "GetScheduleUsers $1,$2", id, tip)
	return users, err
}
