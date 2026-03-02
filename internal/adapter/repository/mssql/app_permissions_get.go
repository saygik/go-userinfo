package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetDomainAccess(upn string) ([]entity.DomainAccess, error) {
	var access []entity.DomainAccess
	_, err := r.db.Select(&access, `
        SELECT domain, access_level
        FROM user_domain_access
        WHERE user_principal_name = $1
    `, upn)
	return access, err
}

func (r *Repository) GetSections(upn string) ([]entity.IdNameDescription, error) {
	var sections []entity.IdNameDescription
	_, err := r.db.Select(&sections, `
        SELECT p.id,
            SUBSTRING(user_permissions.permission_key, 9, LEN(user_permissions.permission_key)-8) as name,
            p.description
        FROM user_permissions
        INNER JOIN permissions p ON user_permissions.permission_key = p.[key]
        WHERE permission_key LIKE 'SECTION:%'
          AND user_principal_name = $1
    `, upn)
	return sections, err
}

func (r *Repository) GetUserRoles(upn string) (roles []entity.IdNameDescription, err error) {
	_, err = r.db.Select(&roles, `
        SELECT r.id, r.name, r.description FROM user_roles ur
        JOIN roles r ON ur.role_id = r.id
        WHERE ur.user_principal_name = $1
    `, upn)
	return roles, err
}

func (r *Repository) GetAppRoles() (roles []entity.IdNameDescription, err error) {
	_, err = r.db.Select(&roles, `
        SELECT id, name, description FROM roles
    `)
	return roles, err
}

func (r *Repository) GetAppSections() (roles []entity.IdNameDescription, err error) {
	_, err = r.db.Select(&roles, `
   SELECT id, SUBSTRING([key], 9, LEN([key])-8) as name, description
        FROM  permissions
        WHERE [key] LIKE 'SECTION:%'
    `)
	return roles, err
}
