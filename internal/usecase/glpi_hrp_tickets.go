package usecase

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetHRPTickets() {

	//	FOR TEST!!!!!!!!!!!!!!!!!!!!
	//	u.matt.SendPostHRP(entity.MattermostHrpPost{})
	// tickets := []entity.GLPI_Ticket{}
	// tickets = append(tickets, entity.GLPI_Ticket{Id: 206238, Content: "Сотрудник: Казаков Юрий Геннадьевич(35407148)"})

	tickets, err := u.glpi.GetHRPTickets()
	_ = tickets
	if err != nil || len(tickets) < 1 {
		return
	}

	redisADUsers, err := u.redis.GetKeyFieldAll("allusers")
	if err != nil {
		return
	}
	sBoxes := ""
	val := ""
	ok := false
	for _, ticket := range tickets {
		finded := false
		channelId, _ := u.glpi.GetGroupMattermostChannel(ticket.GroupId)
		// if err != nil {
		// 	u.log.Error(fmt.Sprintf("Error getting channelId for group %d: %v", ticket.GroupId, err))
		// }

		if len(ticket.Content) < 20 || !strings.Contains(ticket.Content, "Сотрудник:") {
			u.glpi.SetHRPTicket(ticket.Id)
			continue
		}
		sfio := ticket.Content[strings.Index(ticket.Content, "Сотрудник:")+20:]
		endOfFio := strings.Index(sfio, "(")
		if endOfFio > 0 {
			sfio = sfio[:endOfFio]
		} else {
			sfio = "no"
		}
		if len(sfio) < 5 {
			u.glpi.SetHRPTicket(ticket.Id)
			continue
		}

		sDolg := ticket.Content[strings.Index(ticket.Content, "Штатная должность:")+35:]
		endOfDolg := strings.Index(sDolg, "&lt;")
		if endOfDolg > 0 {
			sDolg = sDolg[:endOfDolg]
		} else {
			sDolg = ""
		}
		sMero := ticket.Content[strings.Index(ticket.Content, "Проведено мероприятие:")+44:]
		endOfMero := strings.Index(sMero, "&lt;")
		if endOfDolg > 0 {
			sMero = sMero[:endOfMero]
		} else {
			sMero = ""
		}

		sPred := ticket.Content[strings.Index(ticket.Content, "ОЕ:")+6:]
		endOfPred := strings.Index(sPred, "&lt;")
		if endOfPred > 0 {
			sPred = sPred[:endOfPred]
		} else {
			sPred = ""
		}
		sPred1 := ticket.Content[strings.Index(ticket.Content, "БЕ:")+6:]
		endOfPred = strings.Index(sPred1, "&lt;")
		if endOfPred > 0 {
			sPred1 = sPred1[:endOfPred]
		} else {
			sPred1 = ""
		}
		hrpUser := entity.HRPUser{Id: ticket.Id, FIO: sfio, Dolg: sDolg, Mero: sMero, Company: sPred1 + ", " + sPred}
		upn := ""
		if strings.HasPrefix(ticket.Company, "БЖД > ИВЦ2") || strings.HasPrefix(ticket.Company, "БЖД > ИВЦ3") {
			for _, value := range redisADUsers {
				var user map[string]interface{}
				json.Unmarshal([]byte(value), &user)
				if (fmt.Sprintf("%v", user["displayName"])) == sfio {
					finded = true
					domain := user["domain"].(map[string]interface{})

					upn = fmt.Sprintf("%v", user["userPrincipalName"])
					sBoxes = ""
					sBoxes += `<b>Поиск по ФИО:</b><br>	`
					sBoxes += fmt.Sprintf(`%v найден в домене %v, учетная запись <b>%s</b><br>`, user["displayName"], domain["name"], upn)
					val, ok = user["company"].(string)
					if ok {
						sBoxes += fmt.Sprintf(`<b>Организация:</b> %s<br>`, val)
					}
					val, ok = user["title"].(string)
					if ok {
						sBoxes += fmt.Sprintf(`<b>Должность:</b> %s<br>`, val)
					}
					sBoxes += `рекомендуется направить заявку группе администраторов этого домена для окончательной проверки и отключения`
					u.AddTicketComment(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10, Content: sBoxes})
					softs, err := u.GetUserSoftwares(upn)
					if err == nil && len(softs) > 0 {
						for _, soft := range softs {
							_ = soft
							u.AddTicketComment(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10,
								Content: fmt.Sprintf(`<b>Поиск по учетной записи найденного пользователя:</b><br>
					          %s найден в списке зарегистрированных пользователей системы <b>%s</b><br>
					          рекомендуется направить заявку группе администраторов этой системы (<b>%s</b>) для окончательной проверки и отключения
					 `, upn, soft.Name, soft.GroupName)})

						}
					}
				}
			}
		}
		softs, err := u.GetUserSoftwaresByFio(sfio)
		if err == nil && len(softs) > 0 {
			for _, soft := range softs {
				_ = soft
				finded = true
				adminsChannelId, _ := u.glpi.GetGroupMattermostChannel(int(soft.Groups_id_tech))
				if len(adminsChannelId) > 0 {
					err = u.matt.SendPostHRPSoft(adminsChannelId, hrpUser, soft)
					if err != nil {
						u.log.Error(fmt.Sprintf("Error sending post for  ticket %d to Mattermost channel %s to  ticket soft group %d. Error: %v", ticket.Id, adminsChannelId, soft.Groups_id_tech, err))
					}
				}
				u.AddTicketComment(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10,
					Content: fmt.Sprintf(`<b>Поиск по ФИО:</b><br>
					          <b>%s</b> найден в списке зарегистрированных пользователей системы <b>%s</b><br>
					          рекомендуется направить заявку группе администраторов этой системы (<b>%s</b>) для окончательной проверки и отключения
					 `, sfio, soft.Name, soft.GroupName)})

			}
		}

		sBoxes = ""
		if strings.HasPrefix(ticket.Company, "БЖД > ИВЦ2") {
			boxes, err := u.repo.GetOneDelegateMailBoxes(sfio)
			if err != nil {
				return
			}

			if len(boxes) > 0 {
				finded = true
				sBoxes += fmt.Sprintf(`<b>Поиск в почтовой системе делегированных прав для %s:</b><br>`, sfio)
				for _, box := range boxes {
					sBoxes += fmt.Sprintf(`пользователю делегированы права на почтовый ящик <b>%s</b><br>`, box.Mail)
				}
				sBoxes += "рекомендуется направить заявку группе администраторов почтовой системы вашего региона для окончательной проверки и отключения"
				u.AddTicketComment(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10, Content: sBoxes})

			}
		}
		if strings.HasPrefix(ticket.Company, "БЖД > ИВЦ2") || strings.HasPrefix(ticket.Company, "БЖД > ИВЦ3") {
			if !finded {
				if len(channelId) > 0 {
					err = u.matt.SendPostHRP(channelId, hrpUser)
					if err != nil {
						u.log.Error(fmt.Sprintf("Error sending post for autoresolved ticket %d to Mattermost channel %s to HRP ticket group %d. Error: %v", ticket.Id, channelId, ticket.GroupId, err))
					} else {
						u.log.Info(fmt.Sprintf(`Post for autoresolved ticket %d to Mattermost channel %s to HRP ticket group %d sended`, ticket.Id, channelId, ticket.GroupId))
					}
				}
				u.AddTicketSolution(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10,
					Content: `<b>Поиск не обнаружил пользователя в доменах и системах<br> Заявка закрыта автоматически`})
			}
		}
		u.glpi.SetHRPTicket(ticket.Id)
	}
}
