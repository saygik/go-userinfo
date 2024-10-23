package mattermost

import (
	"context"
	"strconv"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) SendPostHRPSoft(channelId string, hrpUser entity.HRPUser, soft entity.Software) (err error) {
	//post := &model.Post{}
	post := &model.Post{
		ChannelId: channelId,
		Metadata: &model.PostMetadata{
			Priority: &model.PostPriority{
				Priority:     model.NewPointer("important"), // Options: "standard", "important", "urgent"
				RequestedAck: model.NewPointer(true),
			},
		}}

	post.SetProps(map[string]interface{}{
		"attachments": []*model.SlackAttachment{
			{
				AuthorName: "Пользователь найден в системе",
				Text:       "##### " + soft.Name + "\n" + "**ФИО: **" + hrpUser.FIO + ", **Должность: **" + hrpUser.Dolg + ", **Мероприятие: **" + hrpUser.Mero,
				Color:      "#FF2200",
				Title:      "Заявка на отключение учетных данных сотрудника №" + strconv.Itoa(hrpUser.Id),
				TitleLink:  "https://support.rw/front/ticket.form.php?id=" + strconv.Itoa(hrpUser.Id),
				Footer:     hrpUser.Company,
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
