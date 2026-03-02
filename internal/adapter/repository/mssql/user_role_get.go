package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetUserRole(upn string, roleId int) (role entity.IdNameDescription, err error) {
	query := `
		SELECT r.id, r.name, r.description
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_principal_name = $1 AND role_id=$2
	`
	err = r.db.SelectOne(&role, query, upn, roleId)
	if err != nil {
		return role, err
	}

	return role, nil
}

func (r *Repository) GetUserSection(upn string, section string) (role entity.IdNameDescription, err error) {
	query := `
        SELECT p.id,
            SUBSTRING(user_permissions.permission_key, 9, LEN(user_permissions.permission_key)-8) as name,
            p.description
        FROM user_permissions
        INNER JOIN permissions p ON user_permissions.permission_key = p.[key]
        WHERE user_principal_name= $1 AND permission_key = $2
	`
	err = r.db.SelectOne(&role, query, upn, section)
	if err != nil {
		return role, err
	}

	return role, nil
}

func (r *Repository) GetDomainAccessOne(upn string, domain string) (role entity.DomainAccess, err error) {
	query := `
        SELECT domain, access_level
        FROM user_domain_access
        WHERE user_principal_name = $1 AND domain= $2
	`
	err = r.db.SelectOne(&role, query, upn, domain)
	if err != nil {
		return role, err
	}

	return role, nil
}
