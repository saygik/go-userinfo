package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetScheduleUserGroups(id int, tip string) (groups []entity.IdNameType, err error) {
	_, err = r.db.Select(&groups, "GetScheduleUserGroups $1,$2", id, tip)
	return groups, err
}

func (r *Repository) GetScheduleUserAvailableGroups(user string, glpiGroups string) (groups []entity.IdName, err error) {
	_, err = r.db.Select(&groups, "GetScheduleUserAvailableGroups $1,$2", user, glpiGroups)
	return groups, err
}
