package mssql

import (
	"database/sql"
	"errors"

	"github.com/saygik/go-userinfo/internal/entity"
)

type softwareMailTemplateRow struct {
	Subject   string `db:"subject"`
	Body      string `db:"body"`
	UpdatedAt string `db:"updated_at"`
	UpdatedBy string `db:"updated_by"`
}

// GetSoftwareMailTemplate возвращает сохранённый шаблон письма системы.
// Второе значение — признак того, что собственный шаблон задан.
func (r *Repository) GetSoftwareMailTemplate(idSoft int64) (template entity.SoftwareMailTemplate, found bool, err error) {
	var row softwareMailTemplateRow
	err = r.db.SelectOne(&row, `
		SELECT subject,
		       body,
		       COALESCE(CONVERT(varchar(40), updated_at, 127), '') AS updated_at,
		       COALESCE(updated_by, '') AS updated_by
		FROM dbo.SoftwareMailTemplates
		WHERE id_soft = ?`, idSoft)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.SoftwareMailTemplate{}, false, nil
	}
	if err != nil {
		return entity.SoftwareMailTemplate{}, false, err
	}
	return entity.SoftwareMailTemplate{
		Subject:   row.Subject,
		Body:      row.Body,
		IsCustom:  true,
		UpdatedAt: row.UpdatedAt,
		UpdatedBy: row.UpdatedBy,
	}, true, nil
}

// SaveSoftwareMailTemplate создаёт или обновляет шаблон письма системы.
func (r *Repository) SaveSoftwareMailTemplate(idSoft int64, subject string, body string, updatedBy string) error {
	_, err := r.db.Exec(`
		IF EXISTS (SELECT 1 FROM dbo.SoftwareMailTemplates WHERE id_soft = ?)
			UPDATE dbo.SoftwareMailTemplates
			SET subject = ?,
			    body = ?,
			    updated_at = SYSDATETIMEOFFSET(),
			    updated_by = ?
			WHERE id_soft = ?
		ELSE
			INSERT INTO dbo.SoftwareMailTemplates (id_soft, subject, body, updated_by)
			VALUES (?, ?, ?, ?)`,
		idSoft, subject, body, updatedBy, idSoft,
		idSoft, subject, body, updatedBy)
	return err
}

// DeleteSoftwareMailTemplate удаляет собственный шаблон системы.
func (r *Repository) DeleteSoftwareMailTemplate(idSoft int64) error {
	_, err := r.db.Exec(`DELETE FROM dbo.SoftwareMailTemplates WHERE id_soft = ?`, idSoft)
	return err
}
