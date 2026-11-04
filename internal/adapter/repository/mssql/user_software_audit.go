package mssql

import "github.com/saygik/go-userinfo/internal/entity"

// AddUserSoftwareAudit сохраняет запись об успешном изменении связи
// записи пользователя системы. Пустые строки пишутся как NULL.
func (r *Repository) AddUserSoftwareAudit(actor, action string, record entity.SoftUser, softName, details string) error {
	_, err := r.db.Exec(`
        INSERT INTO dbo.UserSoftwareAudit
            (actor_upn, action, user_software_id, id_soft, soft_name,
             user_name, user_login, fio, details)
        VALUES (?, ?, NULLIF(?, 0), NULLIF(?, 0), NULLIF(?, ''),
                NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''), NULLIF(?, ''))
    `, actor, action, record.Id, record.IdSoft, softName,
		record.Name, record.Login, record.Fio, details)
	return err
}

// GetUserSoftwareAudits возвращает последние записи аудита пользователей систем.
func (r *Repository) GetUserSoftwareAudits(limit int) ([]entity.UserSoftwareAudit, error) {
	if limit <= 0 {
		limit = 200
	}
	var rows []entity.UserSoftwareAudit
	_, err := r.db.Select(&rows, `
        SELECT TOP (?)
               id, changed_at, actor_upn, action,
               ISNULL(user_software_id, 0) AS user_software_id,
               ISNULL(id_soft, 0) AS id_soft,
               ISNULL(soft_name, '') AS soft_name,
               ISNULL(user_name, '') AS user_name,
               ISNULL(user_login, '') AS user_login,
               ISNULL(fio, '') AS fio,
               ISNULL(details, '') AS details
        FROM dbo.UserSoftwareAudit
        ORDER BY changed_at DESC, id DESC
    `, limit)
	return rows, err
}
