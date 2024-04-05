package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetUserActivity(login string) (activity []entity.UserActivity, err error) {
	_, err = r.db.Select(&activity, "GetUserLastMonthActivity $1", login)
	return activity, err
}
