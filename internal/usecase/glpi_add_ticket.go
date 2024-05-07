package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

// Create ticket comment...
func (u *UseCase) AddTicket(form entity.NewTicketForm) (int, error) {

	if len(form.User) > 0 {
		token, err := u.glpi.GetUserApiTokenByName(form.User)
		if err != nil {
			return 0, u.Error("ошибка создания  заявки GLPI: у пользователя нет api токена")
		}
		form.Token = token.Name
	}

	if form.UsersIdRequester < 1 {
		return 0, u.Error("ошибка создания  заявки GLPI: отсутствует инициатор")
	}

	formP := entity.NewTicketInputForm{Input: form}
	ticketID, err := u.glpiApi.CreateTicket(formP)
	if ticketID == 0 || err != nil {
		return 0, u.Error("ошибка создания заявки GLPI")
	}
	return ticketID, nil
}
