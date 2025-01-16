package mattermost

import (
	"context"
	"strconv"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) SendPostHRPSoft(channelId string, hrpUser entity.HRPUser, softName string, idTaskCalNotification int) (err error) {

	var actions []*model.PostAction

	if idTaskCalNotification > 0 {
		calAction := model.PostAction{
			Id:    "calendarCompleteTask",
			Name:  "Отменить задачу календаря",
			Style: "primary",
		}
		calActionMap := map[string]any{
			"action": "disable",
			"id":     idTaskCalNotification,
		}

		calIntegration := model.PostActionIntegration{
			URL:     r.integrations.DisableCalendarTaskNotificationApi,
			Context: calActionMap,
		}
		calAction.Integration = &calIntegration
		actions = append(actions, &calAction)

	}
	//******** GLPI Comment add action
	glpiCommentAction := model.PostAction{
		Id:    "glpiCommentAdd",
		Name:  "Комментарий заявки об отключении",
		Style: "success",
	}
	glpiCommentMap := map[string]any{
		"action":  "add",
		"comment": "Комментарий",
		"id":      hrpUser.Id,
		"soft":    softName,
	}
	glpiCommentIntegration := model.PostActionIntegration{
		URL:     r.integrations.AddCommentFromApi,
		Context: glpiCommentMap,
	}
	glpiCommentAction.Integration = &glpiCommentIntegration
	actions = append(actions, &glpiCommentAction)

	//*************************************
	post := &model.Post{
		ChannelId: channelId,
		Metadata: &model.PostMetadata{
			Priority: &model.PostPriority{
				Priority:     model.NewPointer("important"), // Options: "standard", "important", "urgent"
				RequestedAck: model.NewPointer(true),
			},
		}}
	calNotification := ""
	if idTaskCalNotification > 0 {
		calNotification = "\n Срок отключения ещё не наступил, задача напоминания добавлена в календарь вашей группы"
	} else {
		calNotification = ""
	}
	post.SetProps(map[string]interface{}{
		"attachments": []*model.SlackAttachment{
			{
				AuthorName: "Пользователь найден в системе",
				Text: "##### " + softName + "\n" + "*Дата мероприятия: " + hrpUser.Date + "*\n" +
					"**ФИО: **" + hrpUser.FIO + ", **Должность: **" + hrpUser.Dolg + ", **Мероприятие: **" + hrpUser.Mero + calNotification,
				Color:     "#FF2200",
				Title:     "Заявка на отключение учетных данных сотрудника №" + strconv.Itoa(hrpUser.Id),
				TitleLink: "https://support.rw/front/ticket.form.php?id=" + strconv.Itoa(hrpUser.Id),
				Footer:    hrpUser.Company,
				Actions:   actions,
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
