package mssql

import "github.com/saygik/go-userinfo/internal/entity"

// GetAllUserRoles возвращает роли всех пользователей одним запросом.
func (r *Repository) GetAllUserRoles() ([]entity.AdminUserRoleRow, error) {
	var rows []entity.AdminUserRoleRow
	_, err := r.db.Select(&rows, `
        SELECT ur.user_principal_name,
               r.id AS role_id,
               r.name,
               r.description
        FROM user_roles ur
        JOIN roles r ON ur.role_id = r.id
        ORDER BY ur.user_principal_name, r.id
    `)
	return rows, err
}

// GetAllUserDomainAccess возвращает доступы всех пользователей к доменам.
func (r *Repository) GetAllUserDomainAccess() ([]entity.AdminUserDomainAccessRow, error) {
	var rows []entity.AdminUserDomainAccessRow
	_, err := r.db.Select(&rows, `
        SELECT user_principal_name, domain, access_level
        FROM user_domain_access
        ORDER BY user_principal_name, domain
    `)
	return rows, err
}

// GetAllUserSections возвращает все дополнительные разделы пользователей.
func (r *Repository) GetAllUserSections() ([]entity.AdminUserSectionRow, error) {
	var rows []entity.AdminUserSectionRow
	_, err := r.db.Select(&rows, `
        SELECT up.user_principal_name,
               p.id,
               SUBSTRING(up.permission_key, 9, LEN(up.permission_key)-8) as name,
               p.description
        FROM user_permissions up
        INNER JOIN permissions p ON up.permission_key = p.[key]
        WHERE up.permission_key LIKE 'SECTION:%'
        ORDER BY up.user_principal_name, p.id
    `)
	return rows, err
}

// CountUserSoftwares возвращает количество записей пользователей систем.
func (r *Repository) CountUserSoftwares() (int, error) {
	var count int
	err := r.db.SelectOne(&count, `SELECT COUNT(*) FROM dbo.UserSoftwares`)
	return count, err
}

// Ping проверяет доступность базы данных.
func (r *Repository) Ping() error {
	return r.db.Db.Ping()
}
