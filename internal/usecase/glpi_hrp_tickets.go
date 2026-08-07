package usecase

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/agnivade/levenshtein"
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetHRPTickets() {
	if !u.IsAppInitialized() {
		return
	}
	//	FOR TEST!!!!!!!!!!!!!!!!!!!!
	//	u.matt.SendPostHRP(entity.MattermostHrpPost{})
	// tickets := []entity.GLPI_Ticket{}
	// tickets = append(tickets, entity.GLPI_Ticket{Id: 206238, Content: "Сотрудник: Казаков Юрий Геннадьевич(35407148)"})
	//* TEST ***************************************
	tickets, err := u.glpi.GetHRPTickets()
	//	tickets, err := u.glpi.GetHRPTickets()
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
		parts := strings.Split(ticket.Company, ">")
		if len(parts) > 2 {
			parts = parts[:2]
		}
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		region := strings.Join(parts, " > ")
		observeCountTicketsPerRegion(region)

		finded := false
		//_, channelId, _, _ := u.glpi.GetGroupMattermostChannel(ticket.GroupId)

		if len(ticket.Content) < 20 || !strings.Contains(ticket.Content, "Сотрудник:") {
			u.glpi.SetHRPTicket(ticket.Id)
			continue
		}
		sfio := getHRPAttributeFromText(ticket.Content[strings.Index(ticket.Content, "Сотрудник:")+20:], "(")
		if len(sfio) < 5 {
			u.glpi.SetHRPTicket(ticket.Id)
			continue
		}
		sfios := getHRPAttributeFromText(ticket.Content[strings.Index(ticket.Content, "Сотрудник:")+20:], "<")
		sDolg := getHRPAttributeFromText(ticket.Content[strings.Index(ticket.Content, "Штатная должность:")+35:], "<")
		sMero := getHRPAttributeFromText(ticket.Content[strings.Index(ticket.Content, "Проведено мероприятие:")+44:], "<")
		sPred := getHRPAttributeFromText(ticket.Content[strings.Index(ticket.Content, "ОЕ:")+6:], "<")

		sPred1 := getHRPAttributeFromText(ticket.Content[strings.Index(ticket.Content, "БЕ:")+6:], "<")
		sCompany := sPred1
		if len(sPred) > 0 && sPred != "<br>" {
			sCompany = sPred1 + ", " + sPred
		}

		sDate := getHRPAttributeFromText(ticket.Content[strings.Index(ticket.Content, "Дата ограничения действия учетной записи:")+78:], "<")
		dateToNotificate := parseTicketDate(sDate)

		hrpUser := entity.HRPUser{Id: ticket.Id, FIO: sfios, Dolg: sDolg, Mero: sMero, Company: sCompany, Date: sDate}
		hrpUserforHook := entity.HRPUser{Id: ticket.Id, FIO: sfio, Dolg: sDolg, Mero: sMero, Company: sCompany, Date: sDate}
		upn := ""
		u.log.Info("Processing HRP user " + sfio + "...")

		if strings.HasPrefix(ticket.Company, "БЖД > ИВЦ2") {
			u.webClient.AddWebhook(entity.WebhookPayload{Data: hrpUserforHook, WebhookId: strconv.Itoa(ticket.Id)})
		}

		//** Поиск пользователя в доменах//
		for _, value := range redisADUsers {
			var user map[string]any
			json.Unmarshal([]byte(value), &user)
			if user["disabled"] == true {
				continue
			}
			systemUser := fmt.Sprintf("%v", user["displayName"])
			if systemUser == "" {
				continue
			}
			match, log := FuzzyMatchFIO(systemUser, sfio, 2)
			if match {
				//println(log)
				finded = true
				// domain := user["domain"].(map[string]any)
				// domainName := domain["name"].(string)
				domainName, err := getDomainNameFromUser(user)
				if err != nil {
					u.log.Error(fmt.Sprintf("Error getting domain name from user %v: %v", user, err))
					continue
				}
				systemName := fmt.Sprintf("домен %s", domainName)
				logMessage := fmt.Sprintf("%s [%s]", systemUser, log)
				domainAdminsGLPIId := u.ad.GetDomainAdminsGLPI(domainName)
				domainAdminsGLPIName, _, _, _ := u.glpi.GetGroupMattermostChannel(domainAdminsGLPIId)
				u.sendHRPToCalendarAndMattermostChannel(hrpUser, systemName, ticket, sfio, logMessage, dateToNotificate, []entity.SoftwareGroup{{Id: int64(domainAdminsGLPIId)}})

				upn = fmt.Sprintf("%v", user["userPrincipalName"])
				sBoxes = ""
				sBoxes += `<b>Поиск по ФИО:</b><br>	`
				sBoxes += fmt.Sprintf(`%s найден в домене %s, учетная запись <b>%s</b><br>`, hrpUser.FIO, domainName, upn)
				sBoxes += fmt.Sprintf(`%s<br>`, logMessage)
				val, ok = user["company"].(string)
				if ok {
					sBoxes += fmt.Sprintf(`<b>Организация:</b> %s<br>`, val)
				}
				val, ok = user["title"].(string)
				if ok {
					sBoxes += fmt.Sprintf(`<b>Должность:</b> %s<br>`, val)
				}
				sBoxes += fmt.Sprintf(`рекомендуется направить заявку группе администраторов этого домена(%s) для окончательной проверки и отключения`, domainAdminsGLPIName)

				u.AddTicketComment(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10, Content: sBoxes})

				// softs, err := u.GetUserSoftwares(upn)
				// if err == nil && len(softs) > 0 {
				// 	for _, soft := range softs {
				// 		_ = soft
				// 		u.AddTicketComment(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10,
				// 			Content: fmt.Sprintf(`<b>Поиск по учетной записи найденного пользователя:</b><br>
				//           %s найден в списке зарегистрированных пользователей системы <b>%s</b><br>
				//           рекомендуется направить заявку группе администраторов этой системы (<b>%s</b>) для окончательной проверки и отключения
				//  `, upn, soft.Name, soft.GroupName)})

				// 	}
				// }
			}
		}

		//** Поиск пользователя в системах//
		softsUsers, err := u.GetSoftwaresUsers()
		if err == nil && len(softsUsers) > 0 {
			for _, softUser := range softsUsers {
				match, log := FuzzyMatchFIO(softUser.Fio, sfio, 2)
				if match {
					systemName := softUser.SoftName

					logMessage := fmt.Sprintf("%s [%s]", softUser.Fio, log)
					finded = true
					u.sendHRPToCalendarAndMattermostChannel(hrpUser, systemName, ticket, sfio, logMessage, dateToNotificate, softUser.Groups)
					u.AddTicketComment(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10,
						Content: fmt.Sprintf(`<b>Поиск по ФИО:</b><br>
					          %s найден в списке зарегистрированных пользователей системы <b>%s</b><br>
							  %s<br>
					          рекомендуется направить заявку группам администраторов этой системы (<b>%s</b>) для окончательной проверки и отключения
					 `, sfio, systemName, logMessage, softUser.GroupNames)})
				}
			}
		}
		// softs, err := u.GetUserSoftwaresByFio(sfio)
		// if err == nil && len(softs) > 0 {
		// 	for _, soft := range softs {
		// 		_ = soft
		// 		finded = true
		// 		u.sendHRPToCalendarAndMattermostChannel(hrpUser, soft.Name, ticket, sfio, "", dateToNotificate, soft.Groups)

		// 		u.AddTicketComment(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10,
		// 			Content: fmt.Sprintf(`<b>Поиск по ФИО:</b><br>
		// 			          <b>%s</b> найден в списке зарегистрированных пользователей системы <b>%s</b><br>
		// 			          рекомендуется направить заявку группам администраторов этой системы (<b>%s</b>) для окончательной проверки и отключения
		// 			 `, sfio, soft.Name, soft.GroupNames())})

		// 	}
		// }

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
		//|| strings.HasPrefix(ticket.Company, "БЖД > ИВЦ3")
		if strings.HasPrefix(ticket.Company, "БЖД > ИВЦ2") {
			if !finded {
				domainMattermostLogChannels := u.ad.GetMattermostLogChannelsByPrefix(ticket.Company)
				if len(domainMattermostLogChannels) > 0 {
					for _, domainMattermostLogChannel := range domainMattermostLogChannels {
						err = u.matt.SendPostHRP(domainMattermostLogChannel, hrpUser)
						if err != nil {
							u.log.Error(fmt.Sprintf("Error sending post for autoresolved ticket %d to Mattermost channel %s to HRP ticket group %d. Error: %v", ticket.Id, domainMattermostLogChannel, ticket.GroupId, err))
						} else {
							u.log.Info(fmt.Sprintf(`Post for autoresolved ticket %d to Mattermost channel %s to HRP ticket group %d sended`, ticket.Id, domainMattermostLogChannel, ticket.GroupId))
						}
					}
				}
				u.AddTicketSolution(entity.NewCommentForm{ItemId: ticket.Id, ItemType: "Ticket", IsPrivate: true, RequestTypesId: 10,
					Content: `<b>Поиск не обнаружил пользователя в доменах и системах<br> Заявка закрыта автоматически`})
			}
		}
		u.glpi.SetHRPTicket(ticket.Id)
	}
}

func getHRPAttributeFromText(str string, endString string) string {
	sfio := str
	endOfFio := strings.Index(sfio, endString)
	if endOfFio > 0 {
		sfio = sfio[:endOfFio]
	} else {
		sfio = "no"
	}
	return sfio
}

func (u *UseCase) sendHRPToCalendarAndMattermostChannel(
	hrpUser entity.HRPUser,
	softName string,
	ticket entity.GLPI_Ticket,
	sfio string,
	logMessage string,
	dateToNotificate string,
	groups []entity.SoftwareGroup,
) {
	for _, group := range groups {
		adminsName, channelId, calId, _ := u.glpi.GetGroupMattermostChannel(int(group.Id))
		//* TEST ***************************************

		sheduleTaskId := 0
		if dateToNotificate != "" && calId > 0 {
			content := parseHTML(ticket.Content)
			testtask := entity.ScheduleTask{
				Id:             0,
				Idc:            calId,
				Tip:            3,
				Status:         2,
				Title:          fmt.Sprintf(`Отключение в системе %s пользователя %s`, softName, sfio),
				Start:          dateToNotificate,
				End:            "",
				Upn:            "",
				AllDay:         true,
				SendMattermost: true,
				Comment:        "Произвести отключение пользователя по заявке https://support.rw/front/ticket.form.php?id=" + strconv.Itoa(ticket.Id) + "\r\n" + content,
			}
			sheduleTask, err := u.AddScheduleTask(testtask)
			if err == nil {
				sheduleTaskId = sheduleTask.Id
			}
		}
		if len(channelId) > 0 {
			err := u.matt.SendPostHRPSoft(channelId, hrpUser, softName, logMessage, sheduleTaskId)
			if err != nil {
				u.log.Error(fmt.Sprintf("Error sending post for  ticket %d to Mattermost channel %s to  glpi group %s. Error: %v", ticket.Id, channelId, adminsName, err))
			}
		}
	}
}

// normalizeFIO приводит строку к единому виду:
// - нижний регистр,
// - замена ё→е, й→и,
// - удаление ь и ъ,
// - удаление всех небуквенных символов (кроме пробелов),
// - схлопывание множественных пробелов.
func normalizeFIO(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	repl := map[string]string{
		"ё": "е",
		"й": "и",
		"ь": "",
		"ъ": "",
	}
	for old, new := range repl {
		s = strings.ReplaceAll(s, old, new)
	}

	var builder strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			builder.WriteRune(r)
		}
	}
	s = builder.String()

	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

// FuzzyMatchFIO сравнивает две строки ФИО с учётом возможных опечаток.
// Возвращает:
//   - match: true, если суммарное расстояние Левенштейна по частям (Ф, И, О) не превышает threshold;
//   - logMsg: строка с деталями сравнения (полезна для аудита и исправления ошибок).
//
// Если совпадение точное (расстояние 0), logMsg будет пустой.
func FuzzyMatchFIO(a, b string, threshold int) (match bool, logMsg string) {
	normA := normalizeFIO(a)
	normB := normalizeFIO(b)

	// Если одна строка содержит другую (например, "Иванов" и "Иванов Петр") – считаем точным совпадением
	if strings.Contains(normA, normB) || strings.Contains(normB, normA) {
		if a == b {
			return true, "точное совпадение!"
		} else {
			return true, "точное совпадение (с нормализацией!)"
		}
	}

	partsA := strings.Fields(normA)
	partsB := strings.Fields(normB)

	// Если какая-то строка пуста – нет совпадения
	if len(partsA) == 0 || len(partsB) == 0 {
		return false, ""
	}

	// Сравниваем части попарно, начиная с первой.
	// Если количество частей разное, лишние части считаем как их длину.
	minLen := len(partsA)
	if len(partsB) < minLen {
		minLen = len(partsB)
	}

	totalDist := 0

	for i := 0; i < minLen; i++ {
		dist := levenshtein.ComputeDistance(partsA[i], partsB[i])
		totalDist += dist

		if totalDist > threshold {
			return false, "" // превысили порог, лог не нужен
		}
	}

	// Учитываем лишние части
	if len(partsA) > len(partsB) {
		for i := minLen; i < len(partsA); i++ {
			extraLen := len(partsA[i])
			totalDist += extraLen

			if totalDist > threshold {
				return false, ""
			}
		}
	} else if len(partsB) > len(partsA) {
		for i := minLen; i < len(partsB); i++ {
			extraLen := len(partsB[i])
			totalDist += extraLen
			if totalDist > threshold {
				return false, ""
			}
		}
	}

	if totalDist == 0 {
		return true, "точное совпадение!" // точное совпадение
	}

	if totalDist <= threshold {
		// Неполное совпадение – формируем лог
		logMsg = "НЕПОЛНОЕ СОВПАДЕНИЕ!"
		return true, logMsg
	}

	return false, ""
}

func getDomainNameFromUser(user map[string]any) (string, error) {
	// 1. Проверяем user
	if user == nil {
		return "", fmt.Errorf("user is nil")
	}

	// 2. Извлекаем domain безопасно
	domainValue, ok := user["domain"]
	if !ok {
		return "", fmt.Errorf("user[domain] not found")
	}

	domain, ok := domainValue.(map[string]any)
	if !ok {
		return "", fmt.Errorf("user[domain] is not map, got %T", domainValue)
	}

	// 3. Извлекаем domain name
	nameValue, ok := domain["name"]
	if !ok {
		return "", fmt.Errorf("domain[name] not found")
	}

	domainName, ok := nameValue.(string)
	if !ok {
		return "", fmt.Errorf("domain[name] is not string, got %T", nameValue)
	}

	if domainName == "" {
		return "", fmt.Errorf("domain name is empty")
	}

	return domainName, nil
}
