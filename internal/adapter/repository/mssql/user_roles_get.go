package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetUserRoles(userID string) (roles []entity.IdName, err error) {
	_, err = r.db.Select(&roles, "GetUserRoles $1", userID)
	return roles, err
}
