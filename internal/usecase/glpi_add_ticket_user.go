package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

// Create ticket comment...
func (u *UseCase) AddTicketUser(form entity.GLPITicketUserForm) error {

	if len(form.User) > 0 {
		token, err := u.glpi.GetUserApiTokenByName(form.User)
		if err != nil {
			return u.Error("ошибка создания комментария заявки GLPI: у пользователя нет api токена")
		}
		form.Token = token.Name
	}

	formP := entity.GLPITicketUserInputForm{Input: form}
	commentID, err := u.glpiApi.AddTicketUser(formP)
	if commentID == 0 || err != nil {
		return u.Error("ошибка создания пользователя заявки GLPI")
	}
	return nil
}
