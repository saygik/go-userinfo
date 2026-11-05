package mssql

import (
	"strconv"
	"strings"

	"github.com/go-gorp/gorp"
	"github.com/saygik/go-userinfo/internal/entity"
)

// GetSchedules возвращает все календари. Используется админкой и списком
// доступных календарей, поэтому не фильтрует по правам.
func (r *Repository) GetSchedules() (schedules []entity.Schedule, err error) {
	_, err = r.db.Select(&schedules, `
		SELECT id, tip, name, private,
		       ISNULL(caption, '') AS caption,
		       ISNULL(mattermostChannel, '') AS mattermost
		FROM dbo.UserSchedules
		ORDER BY name`)
	return schedules, err
}

// GetScheduleAccessRows возвращает выданные права без раскрытия GLPI-групп
// в пользователей: роль 1 — полный доступ, роль 2 — просмотр и участие
// в графике дежурств.
func (r *Repository) GetScheduleAccessRows(id int) (users []entity.ScheduleAccessUserRow, groups []entity.ScheduleAccessGroupRow, err error) {
	_, err = r.db.Select(&users, `
		SELECT name, tip AS role
		FROM dbo.UserScheduleUsers
		WHERE id_shedule = ?
		ORDER BY name`, id)
	if err != nil {
		return nil, nil, err
	}

	_, err = r.db.Select(&groups, `
		SELECT TRY_CONVERT(INT, name) AS group_id, role
		FROM dbo.UserScheduleUserGroups
		WHERE id_shedule = ? AND tip = 'glpi'
		ORDER BY name`, id)
	return users, groups, err
}

// AddSchedule создаёт календарь и его права в одной транзакции.
func (r *Repository) AddSchedule(form entity.ScheduleForm) (id int, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	err = tx.QueryRow(`
		INSERT INTO dbo.UserSchedules
			(name, tip, mattermostChannel, private, caption)
		OUTPUT INSERTED.id
		VALUES (?, ?, NULLIF(?, ''), ?, NULLIF(?, ''))`,
		form.Name, form.Tip, form.Mattermost,
		form.Private, form.Caption).Scan(&id)
	if err != nil {
		return 0, err
	}

	if err = replaceScheduleAccess(tx, id, form); err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

// UpdateSchedule обновляет календарь и полностью заменяет его права.
func (r *Repository) UpdateSchedule(id int, form entity.ScheduleForm) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		UPDATE dbo.UserSchedules
		SET name = ?,
		    tip = ?,
		    mattermostChannel = NULLIF(?, ''),
		    private = ?,
		    caption = NULLIF(?, '')
		WHERE id = ?`,
		form.Name, form.Tip, form.Mattermost,
		form.Private, form.Caption, id)
	if err != nil {
		return err
	}

	if err = replaceScheduleAccess(tx, id, form); err != nil {
		return err
	}
	return tx.Commit()
}

// DeleteSchedule удаляет календарь вместе с задачами и правами.
func (r *Repository) DeleteSchedule(id int) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
		IF OBJECT_ID(N'dbo.ScheduleTaskMattermostLinks', N'U') IS NOT NULL
			DELETE l
			FROM dbo.ScheduleTaskMattermostLinks l
			INNER JOIN dbo.UserScheduleTasks t ON t.id = l.task_id
			WHERE t.idc = ?`, id)
	if err != nil {
		return err
	}

	if _, err = tx.Exec(`DELETE FROM dbo.UserScheduleTasks WHERE idc = ?`, id); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM dbo.UserScheduleUsers WHERE id_shedule = ?`, id); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM dbo.UserScheduleUserGroups WHERE id_shedule = ?`, id); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM dbo.UserSchedules WHERE id = ?`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func replaceScheduleAccess(tx *gorp.Transaction, id int, form entity.ScheduleForm) error {
	if _, err := tx.Exec(`DELETE FROM dbo.UserScheduleUsers WHERE id_shedule = ?`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM dbo.UserScheduleUserGroups WHERE id_shedule = ?`, id); err != nil {
		return err
	}

	for _, user := range form.Managers.Users {
		if err := insertScheduleUser(tx, id, user, 1); err != nil {
			return err
		}
	}
	for _, user := range form.Viewers.Users {
		if err := insertScheduleUser(tx, id, user, 2); err != nil {
			return err
		}
	}
	for _, group := range form.Managers.Groups {
		if err := insertScheduleGroup(tx, id, group.Id, 1); err != nil {
			return err
		}
	}
	for _, group := range form.Viewers.Groups {
		if err := insertScheduleGroup(tx, id, group.Id, 2); err != nil {
			return err
		}
	}
	return nil
}

func insertScheduleUser(tx *gorp.Transaction, scheduleID int, user string, role int) error {
	_, err := tx.Exec(`
		INSERT INTO dbo.UserScheduleUsers (id_shedule, name, tip)
		VALUES (?, ?, ?)`, scheduleID, user, role)
	return err
}

func insertScheduleGroup(tx *gorp.Transaction, scheduleID int, groupID int, role int) error {
	_, err := tx.Exec(`
		INSERT INTO dbo.UserScheduleUserGroups (id_shedule, tip, name, role)
		VALUES (?, 'glpi', ?, ?)`, scheduleID, strconv.Itoa(groupID), role)
	return err
}

// parseScheduleGroupIDs разбирает строку вида "4,21,34" в список
// положительных идентификаторов. Некорректные значения пропускаются.
func parseScheduleGroupIDs(raw string) []int {
	parts := strings.Split(raw, ",")
	ids := make([]int, 0, len(parts))
	for _, part := range parts {
		id, err := strconv.Atoi(strings.TrimSpace(part))
		if err == nil && id > 0 {
			ids = append(ids, id)
		}
	}
	return ids
}
