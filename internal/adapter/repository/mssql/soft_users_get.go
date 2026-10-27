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
