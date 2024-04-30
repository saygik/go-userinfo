package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetUserGlpiTrackingGroups(userID string) (groups []entity.Id, err error) {
	_, err = r.db.Select(&groups, "GetUserGlpiTrackingGroups $1", userID)
	return groups, err
}
