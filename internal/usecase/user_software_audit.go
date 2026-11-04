package usecase

import (
	"fmt"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

// Действия аудита пользователей систем.
const (
	UserSoftwareAuditCreate = "create"
	UserSoftwareAuditUpdate = "update"
	UserSoftwareAuditDelete = "delete"
)

// auditUserSoftwareChange пишет запись аудита после успешного изменения
// связи пользователя с системой. Ошибка аудита не отменяет изменение.
func (u *UseCase) auditUserSoftwareChange(actor, action string, record entity.SoftUser, details string) {
	if u.repo == nil {
		return
	}
	softName := u.softNameForAudit(record.IdSoft)
	if err := u.repo.AddUserSoftwareAudit(actor, action, record, softName, details); err != nil && u.log != nil {
		u.log.WithError(err).Warn("не удалось записать аудит пользователей систем")
	}
}

// softNameForAudit пытается получить название системы из GLPI.
// Ошибка не мешает записать аудит: название останется пустым.
func (u *UseCase) softNameForAudit(idSoft int) string {
	if u.glpi == nil || idSoft <= 0 {
		return ""
	}
	soft, err := u.glpi.GetSoftware(idSoft)
	if err != nil {
		return ""
	}
	return soft.Name
}

// describeUserSoftwareChanges перечисляет поля, которые изменились.
func describeUserSoftwareChanges(old, updated entity.SoftUser) string {
	fields := make([]string, 0, 8)
	if old.Login != updated.Login {
		fields = append(fields, "логин")
	}
	if old.Fio != updated.Fio {
		fields = append(fields, "ФИО")
	}
	if old.Comment != updated.Comment {
		fields = append(fields, "комментарий")
	}
	if old.Mail != updated.Mail {
		fields = append(fields, "почта")
	}
	if old.GLPITicketID != updated.GLPITicketID {
		fields = append(fields, "заявка GLPI")
	}
	if old.External != updated.External {
		fields = append(fields, "признак внешнего пользователя")
	}
	if !softPermissionsEqual(old.Permissions, updated.Permissions) {
		fields = append(fields, "полномочия")
	}
	if !softAccountStatesEqual(old.AccountState, updated.AccountState) {
		fields = append(fields, "состояние учётной записи")
	}
	if len(fields) == 0 {
		return "без изменений"
	}
	return "изменены: " + strings.Join(fields, ", ")
}

func softPermissionsEqual(left, right []entity.SoftUserPermission) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i].Name != right[i].Name ||
			left[i].EndDate != right[i].EndDate ||
			left[i].Sended != right[i].Sended {
			return false
		}
	}
	return true
}

func softAccountStatesEqual(left, right *entity.SoftUserAccountState) bool {
	if left == nil || right == nil {
		return left == right
	}
	return left.Status == right.Status && left.StatusReason == right.StatusReason
}

// userSoftwareAuditDetails формирует краткое описание изменения.
func userSoftwareAuditDetails(action string, record entity.SoftUser) string {
	switch action {
	case UserSoftwareAuditCreate:
		return fmt.Sprintf("создана запись; полномочий: %d", len(record.Permissions))
	case UserSoftwareAuditDelete:
		return fmt.Sprintf("удалена запись; полномочий: %d", len(record.Permissions))
	default:
		return ""
	}
}

// GetUserSoftwareAudit возвращает последние записи аудита изменений
// связей пользователей с системами.
func (u *UseCase) GetUserSoftwareAudit(limit int) ([]entity.UserSoftwareAudit, error) {
	if limit <= 0 {
		limit = 200
	}
	if limit > 1000 {
		limit = 1000
	}
	rows, err := u.repo.GetUserSoftwareAudits(limit)
	if rows == nil {
		rows = []entity.UserSoftwareAudit{}
	}
	return rows, err
}
