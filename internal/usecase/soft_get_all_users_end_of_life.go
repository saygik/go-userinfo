package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetSoftwareUsersEOL() ([]map[string]interface{}, error) {
	if !u.IsAppInitialized() {
		return nil, u.Error("система не инициализирована")
	}
	permissions, err := u.repo.GetSoftwareUsersEOL()
	if err != nil {
		return nil, u.Error("невозможно получить полномочия с истекающим сроком")
	}
	softUsers := []map[string]interface{}{}

	for _, permission := range permissions {
		soft, err := u.glpi.GetSoftware(permission.IdSoft)
		if err != nil {
			continue
		}

		if len(permission.Mail) > 4 {
			err = u.SendMail(permission.Mail,
				fmt.Sprintf(`
Срок действия полномочий «%s» учетной записи %s пользователя %s для системы %s истекает %s.
Примите меры для перерегистрации в системе или пропустите это письмо, если данная система вас не интересует.`,
					permission.PermissionName, permission.Login, permission.Fio, soft.Name, parseDate(permission.EndDate)))
			if err != nil && u.log != nil {
				u.log.Errorf("не удалось отправить оповещение о полномочиях %d: %v", permission.PermissionId, err)
			}
		}

		for _, group := range soft.Groups {
			if group.GroupCalendar > 0 {
				testtask := entity.ScheduleTask{
					Id:             0,
					Idc:            group.GroupCalendar,
					Tip:            3,
					Status:         2,
					Title:          fmt.Sprintf(`Срок действия полномочий «%s» истёк для пользователя %s`, permission.PermissionName, permission.Fio),
					Start:          permission.EndDate,
					End:            "",
					Upn:            "",
					AllDay:         true,
					SendMattermost: true,
					Comment:        fmt.Sprintf(`Срок действия полномочий «%s» учетной записи %s пользователя %s для системы %s истёк. Произведите отключение.`, permission.PermissionName, permission.Login, permission.Fio, soft.Name),
				}
				_, _ = u.AddScheduleTask(testtask)
			}
		}

		if err = u.repo.SetOneUserSoftwarePermissionSended(permission.PermissionId); err != nil && u.log != nil {
			u.log.Errorf("не удалось отметить оповещение полномочия %d: %v", permission.PermissionId, err)
		}
	}
	return softUsers, nil
}
