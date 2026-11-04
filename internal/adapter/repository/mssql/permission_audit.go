package mssql

import "github.com/saygik/go-userinfo/internal/entity"

// AddPermissionAudit сохраняет запись об успешном изменении прав.
func (r *Repository) AddPermissionAudit(actor, target, action, object, details string) error {
	_, err := r.db.Exec(`
        INSERT INTO dbo.PermissionAudit (actor_upn, target_upn, action, object_name, details)
        VALUES (?, ?, ?, NULLIF(?, ''), NULLIF(?, ''))
    `, actor, target, action, object, details)
	return err
}

// GetPermissionAudits возвращает последние записи аудита изменений прав.
func (r *Repository) GetPermissionAudits(limit int) ([]entity.PermissionAudit, error) {
	if limit <= 0 {
		limit = 200
	}
	var rows []entity.PermissionAudit
	_, err := r.db.Select(&rows, `
        SELECT TOP (?)
               id, changed_at, actor_upn, target_upn, action,
               ISNULL(object_name, '') AS object_name,
               ISNULL(details, '') AS details
        FROM dbo.PermissionAudit
        ORDER BY changed_at DESC, id DESC
    `, limit)
	return rows, err
}
