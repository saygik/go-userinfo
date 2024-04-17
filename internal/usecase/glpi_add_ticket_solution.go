package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

// Create ticket comment...
func (u *UseCase) AddTicketSolution(form entity.NewCommentForm) error {

	if len(form.User) > 0 {
		token, err := u.glpi.GetUserApiTokenByName(form.User)
		if err != nil {
			return u.Error("ошибка создания решения заявки GLPI: у пользователя нет api токена")
		}
		form.Token = token.Name
	}

	formP := entity.NewCommentInputForm{Input: form}
	commentID, err := u.glpiApi.CreateSolution(formP)
	if commentID == 0 || err != nil {
		return u.Error("ошибка создания решения заявки GLPI")
	}
	return nil

}
