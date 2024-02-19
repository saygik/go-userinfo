package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetUserGroups(userID string) (groups []entity.IdName, err error) {
	_, err = r.db.Select(&groups, "GetUserGroups $1", userID)
	return groups, err
}
