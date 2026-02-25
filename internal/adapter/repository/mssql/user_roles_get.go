package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetUserRoles(upn string) (roles []entity.IdNameDescription, err error) {
	_, err = r.db.Select(&roles, `
        SELECT r.id, r.name, r.description FROM user_roles ur
        JOIN roles r ON ur.role_id = r.id
        WHERE ur.user_principal_name = $1
    `, upn)
	return roles, err
}
