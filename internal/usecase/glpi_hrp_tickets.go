package usecase

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetHRPTickets() {

	//FOR TEST!!!!!!!!!!!!!!!!!!!!
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
		if len(ticket.Content) < 20 || !strings.Contains(ticket.Content, "Сотрудник:") {
			u.glpi.SetHRPTicket(ticket.Id)
			continue
		}
		sfio := ticket.Content[strings.Index(ticket.Content, "Сотрудник:")+20:]
		sfio = sfio[:strings.Index(sfio, "(")]
		if len(sfio) < 5 {
			u.glpi.SetHRPTicket(ticket.Id)
			continue
		}
		upn := ""
		for _, value := range redisADUsers {
			var user map[string]interface{}
			json.Unmarshal([]byte(value), &user)
			if (fmt.Sprintf("%v", user["displayName"])) == sfio {
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
		softs, err := u.GetUserSoftwaresByFio(sfio)
		if err == nil && len(softs) > 0 {
			for _, soft := range softs {
				_ = soft
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
				sBoxes += fmt.Sprintf(`<b>Поиск в почтовой системе делегированных прав для %s:</b><br>`, sfio)
				for _, box := range boxes {
					sBoxes += fmt.Sprintf(`пользователю делегированы права на почтовый ящик <b>%s</b><br>`, box.Mail)
				}
				sBoxes += "рекомендуется направить заявку группе администраторов почтовой системы вашего региона для окончательной проверки и отключения"
				u.AddTicketComment(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10, Content: sBoxes})

			}
		}
		if upn == "" {
			u.AddTicketComment(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10,
				Content: `<b>Поиск по ФИО</b><br>не обнаружил пользователя в доменах, подключенных к системе<br> заявка рекомендуется к закрытию(решению)`})

		}
		u.glpi.SetHRPTicket(ticket.Id)
	}
}
