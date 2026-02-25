package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) GetRoles(upn string) ([]string, error) {
	var roles []string
	_, err := r.db.Select(&roles, `
        SELECT r.name FROM user_roles ur
        JOIN roles r ON ur.role_id = r.id
        WHERE ur.user_principal_name = $1
    `, upn)
	return roles, err
}

func (r *Repository) GetDomainAccess(upn string) ([]entity.DomainAccess, error) {
	var access []entity.DomainAccess
	_, err := r.db.Select(&access, `
        SELECT domain, access_level
        FROM user_domain_access
        WHERE user_principal_name = $1
    `, upn)
	return access, err
}

func (r *Repository) GetSections(upn string) ([]string, error) {
	var sections []string
	_, err := r.db.Select(&sections, `
        SELECT SUBSTRING(permission_key, 9, LEN(permission_key))
        FROM user_permissions
        WHERE user_principal_name = $1 AND permission_key LIKE 'SECTION:%'
    `, upn)
	return sections, err
}
