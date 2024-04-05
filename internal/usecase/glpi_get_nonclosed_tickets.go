package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetGLPITicketsNonClosed(user string) ([]entity.GLPI_Ticket, error) {
	ticketsAllowed := []entity.GLPI_Ticket{}
	userRequesterInGLPI, err := u.glpi.GetUserByName(user)
	if err != nil {
		return ticketsAllowed, u.Error("вы не найдены в системе GLPI")
	}
	glpiUserRequesterProfiles, err := u.glpi.GetUserProfiles(userRequesterInGLPI.Id)
	if err == nil {
		if len(glpiUserRequesterProfiles) == 0 {
			return ticketsAllowed, u.Error("ваш профиль не найден в системе GLPI")
		}
		userRequesterInGLPI.Profiles = glpiUserRequesterProfiles
	} else {
		return ticketsAllowed, u.Error("ваш профиль не найден в системе GLPI")
	}
	glpiUserGroups, err := u.glpi.GetUserGroups(userRequesterInGLPI.Id)
	if err == nil {
		userRequesterInGLPI.Groups = glpiUserGroups
	} else {
		userRequesterInGLPI.Groups = []entity.IdName{}
	}
	tickets, err := u.glpi.GetTicketsNonClosed()
	if err != nil {
		return tickets, u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	var ticketOrgs []int
	var ticketUsers []entity.IdNameType
	var ticketGroups []entity.IdNameType

	for _, ticket := range tickets {
		ticket.MyTicket = 0
		if err := json.Unmarshal([]byte(ticket.UsersS), &ticketUsers); err == nil {
			if containsIntInIdNameTypeArray(ticketUsers, 1, userRequesterInGLPI.Id) {
				ticket.UserPresence.Initiator = true
				ticket.MyTicket = 1
			}
			if containsIntInIdNameTypeArray(ticketUsers, 2, userRequesterInGLPI.Id) {
				ticket.UserPresence.Executor = true
				ticket.MyTicket = 1
			}
			if containsIntInIdNameTypeArray(ticketUsers, 3, userRequesterInGLPI.Id) {
				ticket.UserPresence.Observer = true
				ticket.MyTicket = 1
			}
		}
		if err := json.Unmarshal([]byte(ticket.UserGroupsS), &ticketGroups); err == nil {
			if containsIDNameInIdNameTypeArray(ticketGroups, 1, userRequesterInGLPI.Groups) {
				ticket.GroupPresence.Initiator = true
				ticket.MyTicket = 2
			}
			if containsIDNameInIdNameTypeArray(ticketGroups, 2, userRequesterInGLPI.Groups) {
				ticket.GroupPresence.Executor = true
				ticket.MyTicket = 2
			}
			if containsIDNameInIdNameTypeArray(ticketGroups, 3, userRequesterInGLPI.Groups) {
				ticket.GroupPresence.Observer = true
				ticket.MyTicket = 2
			}
		}
		for _, tp := range userRequesterInGLPI.Profiles {
			if tp.Id == 6 {
				if ticket.MyTicket > 0 {
					ticketsAllowed = append(ticketsAllowed, ticket)
					break
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
					ticketsAllowed = append(ticketsAllowed, ticket)
					break

				}
				if containsInt(ticketOrgs, tp.Eid) {
					ticketsAllowed = append(ticketsAllowed, ticket)
					break

				}
			} else {
				if ticket.Eid == tp.Eid {
					ticketsAllowed = append(ticketsAllowed, ticket)
					break
				}
			}
			// if tp.Recursive {
			// }
		}
	}
	return ticketsAllowed, nil
}
