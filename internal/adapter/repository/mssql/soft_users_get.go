package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetSoftwareUsers(softID int) (users []entity.SoftUser, err error) {
	_, err = r.db.Select(&users, "GetSoftwareUsers $1", softID)
	return users, err
}
