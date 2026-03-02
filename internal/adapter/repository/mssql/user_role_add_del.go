package mssql

import "github.com/saygik/go-userinfo/internal/entity"

func (r *Repository) AddUserRole(upn string, roleId int) (*entity.IdNameDescription, error) {

	_, err := r.db.Db.Exec(`
        INSERT INTO user_roles (user_principal_name, role_id)
         VALUES (?, ?)
    `, upn, roleId)
	if err != nil {
		return nil, err
	}
	role, err := r.GetUserRole(upn, roleId)
	if err != nil {
		return &role, err
	}
	return &role, nil
}

func (r *Repository) DelUserGroup(user string, id int) error {
	res, err := r.db.Exec("DelfindUserUserGroup $1,$2", user, id)
	return sqlRowAffectedErrorWrapper(res, err)
}

func (r *Repository) DelUserRole(upn string, roleId int) error {
	res, err := r.db.Exec(`
        DELETE FROM user_roles
        WHERE user_principal_name = ? AND role_id = ?
    `, upn, roleId)
	return sqlRowAffectedErrorWrapper(res, err)
}

func (r *Repository) AddUserSection(upn string, section string) (*entity.IdNameDescription, error) {

	_, err := r.db.Db.Exec(`
        INSERT INTO user_permissions (user_principal_name, permission_key)
         VALUES (?, ?)
    `, upn, section)
	if err != nil {
		return nil, err
	}
	newsection, err := r.GetUserSection(upn, section)
	if err != nil {
		return &newsection, err
	}
	return &newsection, nil
}

func (r *Repository) DelUserSection(upn string, section string) error {
	res, err := r.db.Exec(`
        DELETE FROM user_permissions
        WHERE user_principal_name = ? AND permission_key = ?
    `, upn, section)
	return sqlRowAffectedErrorWrapper(res, err)
}

func (r *Repository) AddUserDomainRole(upn string, domain string, level string) (*entity.DomainAccess, string, error) {
	var count int
	err := r.db.SelectOne(&count, `
        SELECT COUNT(*) FROM user_domain_access
        WHERE user_principal_name = ? AND domain = ?
    `, upn, domain)
	if err != nil {
		return nil, "", err
	}

	operation := "INSERTED"
	if count > 0 {
		operation = "UPDATED"
		// UPDATE
		_, err = r.db.Exec(`
            UPDATE user_domain_access SET access_level = ?
            WHERE user_principal_name = ? AND domain = ?
        `, level, upn, domain)
	} else {
		// INSERT
		_, err = r.db.Exec(`
            INSERT INTO user_domain_access (user_principal_name, domain, access_level)
            VALUES (?, ?, ?)
        `, upn, domain, level)
	}

	if err != nil {
		return nil, "", err
	}
	domainAccess, err := r.GetDomainAccessOne(upn, domain)
	if err != nil {
		return nil, "", err
	}
	// Формируем ответ
	return &domainAccess, operation, nil

}

func (r *Repository) DelUserDomainRole(upn string, domain string) error {
	res, err := r.db.Exec(`
        DELETE FROM user_domain_access
        WHERE user_principal_name = ? AND domain = ?
    `, upn, domain)
	return sqlRowAffectedErrorWrapper(res, err)
}
