package mattermost

import (
	"context"
	"strconv"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) SendPostHRP(channelId string, hrppost entity.HRPUser) (err error) {
	post := &model.Post{
		ChannelId: channelId,
		Metadata: &model.PostMetadata{
			Priority: &model.PostPriority{
				Priority:     model.NewPointer("standard"), // Options: "standard", "important", "urgent"
				RequestedAck: model.NewPointer(true),
			},
		}}

	post.SetProps(map[string]interface{}{

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

	if _, _, err := r.client.CreatePost(context.Background(), post); err != nil {
		return err
	}
	return nil
}
