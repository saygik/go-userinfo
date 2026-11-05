package mssql

import (
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetScheduleUserGroups(id int, tip string) (groups []entity.IdNameType, err error) {
	_, err = r.db.Select(&groups, "GetScheduleUserGroups $1,$2", id, tip)
	return groups, err
}

type scheduleGrant struct {
	IdShedule int `db:"id_shedule"`
	Role      int `db:"role"`
}

// GetScheduleUserAvailableGroups возвращает календари, доступные пользователю.
// type = 2 означает полный доступ, type = 1 — просмотр и возможность быть
// назначенным в события графика дежурств. Для SysAdmin возвращаются все
// календари с полным доступом.
func (r *Repository) GetScheduleUserAvailableGroups(user string, glpiGroups string, isSysAdmin bool) (groups []entity.IdNameType, err error) {
	schedules, err := r.GetSchedules()
	if err != nil {
		return nil, err
	}

	if isSysAdmin {
		result := make([]entity.IdNameType, 0, len(schedules))
		for _, schedule := range schedules {
			result = append(result, entity.IdNameType{Id: schedule.Id, Name: schedule.Name, Type: 2})
		}
		return result, nil
	}

	grants := map[int]int{}
	var userGrants []scheduleGrant
	_, err = r.db.Select(&userGrants, `
		SELECT id_shedule, tip AS role
		FROM dbo.UserScheduleUsers
		WHERE name = ?`, user)
	if err != nil {
		return nil, err
	}
	for _, grant := range userGrants {
		mergeScheduleGrant(grants, grant.IdShedule, grant.Role)
	}

	groupIDs := parseScheduleGroupIDs(glpiGroups)
	if len(groupIDs) > 0 {
		placeholders := make([]string, len(groupIDs))
		args := make([]interface{}, len(groupIDs))
		for i, id := range groupIDs {
			placeholders[i] = "?"
			args[i] = id
		}
		var groupGrants []scheduleGrant
		query := `
			SELECT id_shedule, role
			FROM dbo.UserScheduleUserGroups
			WHERE tip = 'glpi' AND name IN (` + strings.Join(placeholders, ",") + `)`
		_, err = r.db.Select(&groupGrants, query, args...)
		if err != nil {
			return nil, err
		}
		for _, grant := range groupGrants {
			mergeScheduleGrant(grants, grant.IdShedule, grant.Role)
		}
	}

	result := make([]entity.IdNameType, 0, len(schedules))
	for _, schedule := range schedules {
		role, hasAccess := grants[schedule.Id]
		if !hasAccess && schedule.Private {
			continue
		}
		accessType := 1
		if role == 1 {
			accessType = 2
		}
		result = append(result, entity.IdNameType{
			Id:   schedule.Id,
			Name: schedule.Name,
			Type: accessType,
		})
	}
	return result, nil
}

func mergeScheduleGrant(grants map[int]int, scheduleID, role int) {
	if scheduleID <= 0 {
		return
	}
	// Роль 1 (полный доступ) важнее роли 2 (просмотр).
	if current, exists := grants[scheduleID]; exists && current == 1 {
		return
	}
	if role != 1 && role != 2 {
		role = 2
	}
	grants[scheduleID] = role
}
