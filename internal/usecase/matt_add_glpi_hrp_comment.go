package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) AddGLPI_HRPTicketCommentFromMattermost(form entity.MattermostInteractiveMessageRequestForm) (string, error) {
	strMattUser, loginMattUser, err := u.GetMattermostUserById(form.UserId)
	if err != nil {
		return "", u.Error("Комментарий не добавлен. Невозможно получить автора комментария")

	}
	content := fmt.Sprintf(`Комментарий пользователя Mattermost <b>%s</b>
пользователь отключен в системе <b>%s</b>`, strMattUser, form.Context.Soft)
	var commentForm entity.NewCommentForm
	commentForm.ItemId = form.Context.Id
	commentForm.RequestTypesId = 11
	commentForm.ItemType = "Ticket"
	commentForm.User = "glpi_find_bot@local"
	commentForm.Content = content

	err = u.AddTicketComment(commentForm)
	if err != nil {
		return "", u.Error("Комментарий не добавлен. Ошибка: " + err.Error())
	}
	return "Комментарий об отключении пользователя добавлен в заявку" + " от пользователя @" + loginMattUser, nil
}
