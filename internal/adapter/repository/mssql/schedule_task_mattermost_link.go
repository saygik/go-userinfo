package mssql

import (
	"database/sql"

	"github.com/saygik/go-userinfo/internal/entity"
)

// SaveScheduleTaskMattermostLink сохраняет связь задачи календаря с постом
// Mattermost. Повторная вставка для той же задачи игнорируется.
func (r *Repository) SaveScheduleTaskMattermostLink(link entity.ScheduleTaskMattermostLink) error {
	_, err := r.db.Exec(`
		IF NOT EXISTS (
			SELECT 1 FROM dbo.ScheduleTaskMattermostLinks WHERE task_id = ?
		)
		INSERT INTO dbo.ScheduleTaskMattermostLinks
			(task_id, channel_id, post_id, source_type, source_id)
		VALUES (?, ?, ?, ?, NULLIF(?, 0))
	`, link.TaskID, link.TaskID, link.ChannelID, link.PostID, link.SourceType, link.SourceID)
	return err
}

// GetScheduleTaskMattermostLink возвращает связь задачи с постом.
// Если связи нет, возвращает пустую структуру без ошибки.
func (r *Repository) GetScheduleTaskMattermostLink(taskID int64) (entity.ScheduleTaskMattermostLink, error) {
	var link entity.ScheduleTaskMattermostLink
	err := r.db.SelectOne(&link, `
		SELECT task_id,
		       channel_id,
		       post_id,
		       source_type,
		       ISNULL(source_id, 0) AS source_id,
		       reminder_sent_at,
		       closed_at,
		       ISNULL(closed_by, '') AS closed_by,
		       ISNULL(close_reason, '') AS close_reason
		FROM dbo.ScheduleTaskMattermostLinks
		WHERE task_id = ?
	`, taskID)
	if err == sql.ErrNoRows {
		return link, nil
	}
	return link, err
}

// MarkScheduleTaskMattermostReminderSent фиксирует факт отправки напоминания.
func (r *Repository) MarkScheduleTaskMattermostReminderSent(taskID int64) error {
	_, err := r.db.Exec(`
		UPDATE dbo.ScheduleTaskMattermostLinks
		SET reminder_sent_at = SYSDATETIMEOFFSET()
		WHERE task_id = ?
		  AND reminder_sent_at IS NULL
	`, taskID)
	return err
}

// MarkScheduleTaskMattermostClosed фиксирует закрытие задачи.
func (r *Repository) MarkScheduleTaskMattermostClosed(taskID int64, actor, reason string) error {
	_, err := r.db.Exec(`
		UPDATE dbo.ScheduleTaskMattermostLinks
		SET closed_at = COALESCE(closed_at, SYSDATETIMEOFFSET()),
		    closed_by = COALESCE(NULLIF(?, ''), closed_by),
		    close_reason = COALESCE(NULLIF(?, ''), close_reason)
		WHERE task_id = ?
	`, actor, reason, taskID)
	return err
}
