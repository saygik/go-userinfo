package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetUserSoftwares(user string) (softwares []entity.IdName, err error) {

	_, err = r.db.Select(&softwares, "GetUserSoftwares $1", user)
	return softwares, err
}

func (r *Repository) GetUserSoftwaresByFio(user string) ([]int64, error) {
	var softwares []int64
	_, err := r.db.Select(&softwares, "GetUserSoftwaresByFio $1", user)
	return softwares, err
}
