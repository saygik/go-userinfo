package usecase

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetGLPITicket(id string, user string) (entity.GLPI_Ticket, error) {
	if id == "" {
		return entity.GLPI_Ticket{}, u.Error("неверное ID заявки")
	}
	userRequesterInGLPI, err := u.glpi.GetUserByName(user)
	if err != nil {
		return entity.GLPI_Ticket{}, u.Error("вы не найдены в системе GLPI")
	}
	glpiUserRequesterProfiles, err := u.glpi.GetUserProfiles(userRequesterInGLPI.Id)
	if err == nil {
		if len(glpiUserRequesterProfiles) == 0 {
			return entity.GLPI_Ticket{}, u.Error("ваш профиль не найден в системе GLPI")
		}
		userRequesterInGLPI.Profiles = glpiUserRequesterProfiles
	} else {
		return entity.GLPI_Ticket{}, u.Error("ваш профиль не найден в системе GLPI")
	}
	glpiUserGroups, err := u.glpi.GetUserGroups(userRequesterInGLPI.Id)
	if err == nil {
		userRequesterInGLPI.Groups = glpiUserGroups
	} else {
		userRequesterInGLPI.Groups = []entity.IdName{}
	}

	ticket, err := u.glpi.GetTicket(id)
	if err != nil {
		return entity.GLPI_Ticket{}, u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	works, _ := u.glpi.GetTicketWorks(id)
	ticket.Works = works
	if ticket.RecipientId > 0 {
		recipient, err := u.glpi.GetUserById(ticket.RecipientId)
		if err == nil {
			adUser := u.GetUserADPropertysShort(recipient.Name)
			if adUser["findedInAD"] == false {
				adUser["displayName"] = recipient.Fio
				adUser["mail"] = recipient.Email
			}
			ticket.Recipient = adUser
		}

	}

	var ADiniciators []map[string]interface{}
	var ADexetutors []map[string]interface{}
	var ADobservers []map[string]interface{}

	//	if err := json.Unmarshal([]byte(ticket.UsersS), &ticketUsers); err == nil {
	users, err := u.glpi.GetTicketUsers(id)
	if err == nil {
		//		ticket.Users = ticketUsers
		for _, user := range users {
			adUser := u.GetUserADPropertysShort(user.Name)
			if adUser["findedInAD"] == false {
				adUser["displayName"] = user.Fio
				adUser["mail"] = user.Email
			}
			switch user.Type {
			case 1:
				if userRequesterInGLPI.Id == user.Id {
					ticket.UserPresence.Initiator = true
					ticket.MyTicket = 1
				}
				ADiniciators = append(ADiniciators, adUser)
			case 2:
				if userRequesterInGLPI.Id == user.Id {
					ticket.UserPresence.Executor = true
					ticket.MyTicket = 1
				}
				ADexetutors = append(ADexetutors, adUser)
			case 3:
				if userRequesterInGLPI.Id == user.Id {
					ticket.UserPresence.Observer = true
					ticket.MyTicket = 1
				}
				ADobservers = append(ADobservers, adUser)
			}
		}
		ticket.Users.Initiators = ADiniciators
		ticket.Users.Executors = ADexetutors
		ticket.Users.Observers = ADobservers
	}
	ticket.UsersS = ""
	groups, err := u.glpi.GetTicketGroupForCurrentUser(id, userRequesterInGLPI.Id)

	Giniciators := []entity.GLPIGroup{}
	Gexetutors := []entity.GLPIGroup{}
	Gobservers := []entity.GLPIGroup{}
	if err == nil {
		for _, group := range groups {
			switch group.Type {
			case 1:
				if group.Presence {
					ticket.GroupPresence.Initiator = true
					ticket.MyTicket = 2
				}
				Giniciators = append(Giniciators, group)
			case 2:
				if group.Presence {
					ticket.GroupPresence.Executor = true
					ticket.MyTicket = 2
				}
				Gexetutors = append(Gexetutors, group)
			case 3:
				if group.Presence {
					ticket.GroupPresence.Observer = true
					ticket.MyTicket = 2
				}
				Gobservers = append(Gobservers, group)
			}
		}
		ticket.Groups.Initiators = Giniciators
		ticket.Groups.Executors = Gexetutors
		ticket.Groups.Observers = Gobservers

	}

	problems, err := u.glpi.GetTicketProblems(id)
	if err == nil {
		ticket.Problems = problems
	}

	ticket.UserGroupsS = ""
	var ticketOrgs []int
	for _, tp := range userRequesterInGLPI.Profiles {
		if tp.Id == 6 {
			if ticket.MyTicket > 0 {
				return ticket, nil
			}
		}
		if tp.Id != 3 && tp.Id != 4 && tp.Id != 5 && tp.Id != 15 && tp.Id != 7 {
			continue
		}
		if tp.Recursive {
			if err := json.Unmarshal([]byte(ticket.Orgs), &ticketOrgs); err != nil {
				continue
			}
			if ticket.Eid == tp.Eid {
				return ticket, nil

			}
			if containsInt(ticketOrgs, tp.Eid) {
				return ticket, nil

			}
		} else {
			if ticket.Eid == tp.Eid {
				return ticket, nil

			}
		}
	}
	return entity.GLPI_Ticket{}, u.Error("у вас нет прав на эту заявку GLPI")
}

func (u *UseCase) GetGLPITicketSimple(id string) (entity.GLPI_Ticket, error) {
	if id == "" {
		return entity.GLPI_Ticket{}, u.Error("неверное ID заявки")
	}

	ticket, err := u.glpi.GetTicket(id)
	if err != nil {
		return entity.GLPI_Ticket{}, u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}

	return ticket, nil
}

func (u *UseCase) GetGLPITicketReport(id string, user string) (*entity.GLPI_Ticket_Report, error) {
	ti := entity.GLPI_Ticket_Report{}
	if id == "" || id == "0" {
		return &ti, u.Error("неверное ID заявки")
	}
	userRequesterInGLPI, err := u.glpi.GetUserByName(user)
	_ = userRequesterInGLPI
	if err != nil {
		return &ti, u.Error("вы не найдены в системе GLPI")
	}
	glpiUserRequesterProfiles, err := u.glpi.GetUserProfiles(userRequesterInGLPI.Id)
	if err == nil {
		if len(glpiUserRequesterProfiles) == 0 {
			return &ti, u.Error("ваш профиль не найден в системе GLPI")
		}
		userRequesterInGLPI.Profiles = glpiUserRequesterProfiles
	} else {
		return &ti, u.Error("ваш профиль не найден в системе GLPI")
	}
	glpiUserGroups, err := u.glpi.GetUserGroups(userRequesterInGLPI.Id)
	if err == nil {
		userRequesterInGLPI.Groups = glpiUserGroups
	} else {
		userRequesterInGLPI.Groups = []entity.IdName{}
	}

	ticket, err := u.glpi.GetTicketReport(id)
	if err != nil {
		return &ti, u.Error("ошибка доступа к базе данных GLPI")
	}

	ticket.Date = parseGlpiDate(ticket.Date)
	ticket.SolveDate = parseGlpiDate(ticket.SolveDate)
	ticket.Closedate = parseGlpiDate(ticket.Closedate)
	ticket.DateCreation = parseGlpiDate(ticket.DateCreation)
	ticket.DateMod = parseGlpiDate(ticket.DateMod)

	var ADiniciators []entity.SimpleUser
	var ADexetutors []entity.SimpleUser
	var ADobservers []entity.SimpleUser

	//	if err := json.Unmarshal([]byte(ticket.UsersS), &ticketUsers); err == nil {
	users, err := u.glpi.GetTicketUsers(id)
	if err == nil {
		//		ticket.Users = ticketUsers
		for _, user := range users {
			adUserPtr, err := u.GetUserADPropertysSimple(user.Name)
			adUser := *adUserPtr
			if err != nil {
				adUser.Name = user.Fio
				adUser.Mail = user.Email
			}
			switch user.Type {
			case 1:
				ADiniciators = append(ADiniciators, adUser)
				if userRequesterInGLPI.Id == user.Id {
					ticket.MyTicket = 1
				}
			case 2:
				ADexetutors = append(ADexetutors, adUser)
				if userRequesterInGLPI.Id == user.Id {
					ticket.MyTicket = 1
				}

			case 3:
				ADobservers = append(ADobservers, adUser)
				if userRequesterInGLPI.Id == user.Id {
					ticket.MyTicket = 1
				}
			}
		}
		ticket.Users.Initiators = ADiniciators
		ticket.Users.Executors = ADexetutors
		ticket.Users.Observers = ADobservers
	}
	works, _ := u.glpi.GetTicketComments(id)

	for i, work := range works {
		adUserPtr, err := u.GetUserADPropertysSimple(work.Name)
		adUser := *adUserPtr
		if err != nil {
			adUser.Name = work.Author
			adUser.Department = ""
			adUser.Title = ""
		}
		works[i].AuthorProps = adUser
		works[i].DateCreation = parseGlpiDate(work.DateCreation)
	}
	ticket.Comments = works
	works, _ = u.glpi.GetTicketSolutions(id)
	for i, work := range works {
		adUserPtr, err := u.GetUserADPropertysSimple(work.Name)
		adUser := *adUserPtr
		if err != nil {
			adUser.Name = work.Author
			adUser.Department = ""
			adUser.Title = ""
		}
		works[i].AuthorProps = adUser
		works[i].DateCreation = parseGlpiDate(work.DateCreation)
	}
	ticket.Solutions = works
	group, err := u.glpi.GetTicketGroupExecutors(id)
	if err == nil {
		ticket.ExecutorsGroup = group.Name
		users, err := u.glpi.GetGroupUsersFio(group.Id)
		if err == nil {
			var ExecutorsUsers []entity.SimpleUser
			for _, user := range users {
				adUserPtr, err := u.GetUserADPropertysSimple(user.Name)
				adUser := *adUserPtr
				if err != nil {
					adUser.Name = user.Fio
					adUser.Department = ""
					adUser.Title = ""
				}
				ExecutorsUsers = append(ExecutorsUsers, adUser)
			}
			ticket.ExecutorsGroupUsers = ExecutorsUsers
		}

	}
	ticket.Objects.NetworkEquipment, _ = u.glpi.GetTicketNetworkEquipment(id)
	ticket.Objects.Servers, _ = u.glpi.GetTicketServers(id)
	ticket.Objects.Softwares, _ = u.glpi.GetTicketSoft(id)
	access := isAccessToTicket(userRequesterInGLPI, ticket)
	if !access {
		return &ti, u.Error("ваши профили GLPI не дают доступа к этой заявке")
	}
	return &ticket, nil
}

func isAccessToTicket(user entity.GLPIUser, ticket entity.GLPI_Ticket_Report) bool {
	var ticketOrgs []int
	technicalProfiles := []int{6, 12, 14, 16, 17}
	administratorProfiles := []int{2, 3, 4, 5, 7, 11, 15}
	if len(ticket.FailCategory) > 0 {
		return true
	}
	for _, tp := range user.Profiles {
		if slices.Contains(technicalProfiles, tp.Id) {
			if ticket.MyTicket > 0 {
				return true
			}
		}
		if !slices.Contains(administratorProfiles, tp.Id) {
			continue
		}
		if tp.Recursive {
			if err := json.Unmarshal([]byte(ticket.Orgs), &ticketOrgs); err != nil {
				continue
			}
			if ticket.Eid == tp.Eid {
				return true

			}
			if containsInt(ticketOrgs, tp.Eid) {
				return true

			}
		} else {
			if ticket.Eid == tp.Eid {
				return true

			}
		}
	}
	return false
}
