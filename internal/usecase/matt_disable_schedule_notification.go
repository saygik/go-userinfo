package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) DisableSheduleTaskNotificationFromMattermost(form entity.MattermostInteractiveMessageRequestForm) (string, error) {
	_, loginMattUser, err := u.GetMattermostUserById(form.UserId)
	if err != nil {
		return "", u.Error("Задача календаря не отменена. Невозможно получить автора действия")
	}

	err = u.repo.UpdateScheduleTaskDisableMattermost(form.Context.Id)
	if err != nil {
		return "", u.Error("Задача календаря не отменена. Ошибка: " + err.Error())
	}
	return "Задача календаря отменена. Оповещение отключено" + " пользователем @" + loginMattUser, nil
}
