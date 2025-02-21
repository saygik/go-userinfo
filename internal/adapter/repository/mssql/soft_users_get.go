package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetSoftwareUsers(softID int) (users []entity.SoftUser, err error) {
	_, err = r.db.Select(&users, "GetSoftwareUsers_2 $1", softID)
	return users, err
}

func (r *Repository) GetSoftwareUsersEOL() (users []entity.SoftUser, err error) {
	_, err = r.db.Select(&users, "GetSoftwareUsersEOL")
	return users, err
}
