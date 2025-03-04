package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetSoftwareUsersEOL() ([]map[string]interface{}, error) {

	users, err := u.repo.GetSoftwareUsersEOL()
	if err != nil {
		return nil, u.Error("невозможно получить системы из GLPI")
	}
	softUsers := []map[string]interface{}{}

	for _, user := range users {

		soft, err := u.glpi.GetSoftware(user.IdSoft)
		if err != nil {
			continue
		}
		if len(user.Mail) > 4 {
			err = u.SendMail(user.Mail,
				fmt.Sprintf(`
Срок действия учетной записи % s пользователя %s для системы %s истекает %s.
Примите меры для перерегистрации в системе или пропустите это письмо, если данная система вас не интересует.`,
					user.Login, user.Fio, soft.Name, parseDate(user.EndDate)))
		}
		if soft.GroupCalendar > 0 {
			testtask := entity.ScheduleTask{
				Id:             0,
				Idc:            soft.GroupCalendar,
				Tip:            3,
				Status:         2,
				Title:          fmt.Sprintf(`Срок действия учетной записи истёк для пользователя %s`, user.Fio),
				Start:          user.EndDate,
				End:            "",
				Upn:            "",
				AllDay:         true,
				SendMattermost: true,
				Comment:        fmt.Sprintf(`Срок действия учетной записи % s пользователя %s для системы %s истёк. Произведите отключение.`, user.Login, user.Fio, soft.Name),
			}
			_, err := u.AddScheduleTask(testtask)
			if err == nil {
				u.repo.SetOneUserSoftwareSendedToCalendar(user.Id)
			}
		}
	}
	return softUsers, nil
}
