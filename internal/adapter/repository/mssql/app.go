package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetUserResourceAccess(resouceID string, userID string) (access int, err error) {
	err = r.db.QueryRow("GetUserResourceAccess $1,$2", resouceID, userID).Scan(&access)
	return access, err
}

func (r *Repository) GetCurrentUserResources(userID string) (groups []entity.AppResource, err error) {
	_, err = r.db.Select(&groups, "GetUserResources $1", userID)
	return groups, err
}
