package mssql

import "github.com/saygik/go-userinfo/internal/entity"

const softUserSelectColumns = `
	SELECT u.id,
	       u.id_soft,
	       u.user_name,
	       COALESCE(u.user_login, '') AS user_login,
	       COALESCE(u.user_comment, '') AS user_comment,
	       COALESCE(u.fio, '') AS fio,
	       u.[external],
	       COALESCE(u.mail, '') AS mail,
	       COALESCE(u.editor, '') AS editor,
	       COALESCE(CONVERT(varchar(40), u.last_changes, 127), '') AS last_changes
	FROM dbo.UserSoftwares u`

const softPermissionSelectColumns = `
	SELECT p.id,
	       p.id_user_software,
	       p.name,
	       COALESCE(CONVERT(varchar(40), p.enddate, 127), '') AS enddate,
	       p.sended
	FROM dbo.UserSoftwaresPermissions p`

const softAccountStateSelectColumns = `
	SELECT a.id,
	       a.id_user_software,
	       a.status,
	       COALESCE(CONVERT(varchar(40), a.status_changed_at, 127), '') AS status_changed_at,
	       COALESCE(a.status_changed_by, '') AS status_changed_by,
	       COALESCE(a.status_reason, '') AS status_reason,
	       COALESCE(a.source, '') AS source,
	       COALESCE(CONVERT(varchar(40), a.last_sync_at, 127), '') AS last_sync_at
	FROM dbo.UserSoftwaresAccountStates a`

func defaultSoftAccountState() *entity.SoftUserAccountState {
	return &entity.SoftUserAccountState{
		Status: entity.SoftUserAccountStatusActive,
		Source: "default",
	}
}

func (r *Repository) GetSoftwareUsers(softID int) (users []entity.SoftUser, err error) {
	_, err = r.db.Select(&users, softUserSelectColumns+`
		WHERE u.id_soft = $1
		ORDER BY u.fio, u.user_name`, softID)
	if err != nil {
		return users, err
	}

	permissions, err := r.getPermissionsBySoftID(softID)
	if err != nil {
		return users, err
	}
	attachSoftPermissions(users, permissions)

	states, err := r.getAccountStatesBySoftID(softID)
	if err != nil {
		return users, err
	}
	attachSoftAccountStates(users, states)
	return users, nil
}

func (r *Repository) GetSoftwaresUsers() (users []entity.SoftUser, err error) {
	_, err = r.db.Select(&users, softUserSelectColumns+`
		ORDER BY u.fio, u.user_name`)
	if err != nil {
		return users, err
	}

	permissions, err := r.getAllPermissions()
	if err != nil {
		return users, err
	}
	attachSoftPermissions(users, permissions)

	states, err := r.getAllAccountStates()
	if err != nil {
		return users, err
	}
	attachSoftAccountStates(users, states)
	return users, nil
}

func (r *Repository) GetSoftwareUsersEOL() (users []entity.SoftUserPermissionEOL, err error) {
	_, err = r.db.Select(&users, `
		SELECT p.id AS permission_id,
		       p.id_user_software,
		       u.id_soft,
		       u.user_name,
		       COALESCE(u.user_login, '') AS user_login,
		       COALESCE(u.fio, '') AS fio,
		       COALESCE(u.mail, '') AS mail,
		       p.name AS permission_name,
		       COALESCE(CONVERT(varchar(40), p.enddate, 127), '') AS enddate
		FROM dbo.UserSoftwaresPermissions p
		INNER JOIN dbo.UserSoftwares u ON u.id = p.id_user_software
		WHERE p.sended = 0
		  AND p.enddate IS NOT NULL
		  AND p.enddate > CONVERT(datetimeoffset(7), '1900-01-02T00:00:00+00:00')
		  AND p.enddate <= DATEADD(day, 30, SYSDATETIMEOFFSET())
		ORDER BY p.enddate, p.id`)
	return users, err
}

func (r *Repository) getPermissionsBySoftID(softID int) (map[int64][]entity.SoftUserPermission, error) {
	var permissions []entity.SoftUserPermission
	_, err := r.db.Select(&permissions, softPermissionSelectColumns+`
		INNER JOIN dbo.UserSoftwares u ON u.id = p.id_user_software
		WHERE u.id_soft = $1
		ORDER BY p.id`, softID)
	if err != nil {
		return nil, err
	}
	return groupSoftPermissions(permissions), nil
}

func (r *Repository) getAllPermissions() (map[int64][]entity.SoftUserPermission, error) {
	var permissions []entity.SoftUserPermission
	_, err := r.db.Select(&permissions, softPermissionSelectColumns+`
		ORDER BY p.id`)
	if err != nil {
		return nil, err
	}
	return groupSoftPermissions(permissions), nil
}

func groupSoftPermissions(permissions []entity.SoftUserPermission) map[int64][]entity.SoftUserPermission {
	result := make(map[int64][]entity.SoftUserPermission)
	for _, permission := range permissions {
		result[permission.IdUserSoft] = append(result[permission.IdUserSoft], permission)
	}
	return result
}

func attachSoftPermissions(users []entity.SoftUser, permissions map[int64][]entity.SoftUserPermission) {
	for i := range users {
		users[i].Permissions = permissions[users[i].Id]
		if users[i].Permissions == nil {
			users[i].Permissions = []entity.SoftUserPermission{}
		}
	}
}

func (r *Repository) getAccountStatesBySoftID(softID int) (map[int64]*entity.SoftUserAccountState, error) {
	var states []entity.SoftUserAccountState
	_, err := r.db.Select(&states, softAccountStateSelectColumns+`
		INNER JOIN dbo.UserSoftwares u ON u.id = a.id_user_software
		WHERE u.id_soft = $1`, softID)
	if err != nil {
		return nil, err
	}
	return groupSoftAccountStates(states), nil
}

func (r *Repository) getAllAccountStates() (map[int64]*entity.SoftUserAccountState, error) {
	var states []entity.SoftUserAccountState
	_, err := r.db.Select(&states, softAccountStateSelectColumns)
	if err != nil {
		return nil, err
	}
	return groupSoftAccountStates(states), nil
}

func groupSoftAccountStates(states []entity.SoftUserAccountState) map[int64]*entity.SoftUserAccountState {
	result := make(map[int64]*entity.SoftUserAccountState, len(states))
	for i := range states {
		result[states[i].IdUserSoft] = &states[i]
	}
	return result
}

func attachSoftAccountStates(users []entity.SoftUser, states map[int64]*entity.SoftUserAccountState) {
	for i := range users {
		if state := states[users[i].Id]; state != nil {
			users[i].AccountState = state
		} else {
			users[i].AccountState = defaultSoftAccountState()
		}
	}
}
