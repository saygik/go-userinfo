package mattermost

import (
	"strconv"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) SendPostHRP(channelId string, hrppost entity.MattermostHrpPost) (err error) {
	post := &model.Post{}

	post.ChannelId = channelId

	post.SetProps(map[string]interface{}{
		// "priority": map[string]interface{}{
		// 	"priority": "important", // Options: "standard", "important", "urgent"
		// },
		"attachments": []*model.SlackAttachment{
			{

				AuthorName: "Решена автоматически роботом",
				Text:       "**ФИО: **" + hrppost.FIO + ", **Должность: **" + hrppost.Dolg + ", **Мероприятие: **" + hrppost.Mero,
				Color:      "#FFA500",
				Title:      "Заявка на отключение учетных данных сотрудника №" + strconv.Itoa(hrppost.Id),
				TitleLink:  "https://support.rw/front/ticket.form.php?id=" + strconv.Itoa(hrppost.Id),
				Footer:     hrppost.Company,
				// Fields: []*model.SlackAttachmentField{
				// 	{
				// 		Title: "ФИО",
				// 		Value: hrppost.FIO,
				// 	},
				// 	{
				// 		Title: "должность",
				// 		Value: hrppost.Dolg,
				// 	},
				// },
			},
		},
	})

	if _, _, err := r.client.CreatePost(post); err != nil {
		return err
	}
	return nil
}
