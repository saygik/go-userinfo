package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetGLPIProblem(id string, user string) (entity.GLPI_Problem, error) {
	if id == "" {
		return entity.GLPI_Problem{}, u.Error("неверное ID заявки")
	}

	userRequesterInGLPI, err := u.glpi.GetUserByName(user)
	if err != nil {
		return entity.GLPI_Problem{}, u.Error("вы не найдены в системе GLPI")
	}
	glpiUserRequesterProfiles, err := u.glpi.GetUserProfiles(userRequesterInGLPI.Id)
	if err == nil {
		if len(glpiUserRequesterProfiles) == 0 {
			return entity.GLPI_Problem{}, u.Error("ваш профиль не найден в системе GLPI")
		}
		userRequesterInGLPI.Profiles = glpiUserRequesterProfiles
	} else {
		return entity.GLPI_Problem{}, u.Error("ваш профиль не найден в системе GLPI")
	}
	glpiUserGroups, err := u.glpi.GetUserGroups(userRequesterInGLPI.Id)
	if err == nil {
		userRequesterInGLPI.Groups = glpiUserGroups
	} else {
		userRequesterInGLPI.Groups = []entity.IdName{}
	}

	problem, err := u.glpi.GetProblem(id)
	if err != nil {
		return entity.GLPI_Problem{}, u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	works, _ := u.glpi.GetProblemWorks(id)
	problem.Works = works
	if problem.RecipientId > 0 {
		recipient, err := u.glpi.GetUserById(problem.RecipientId)
		if err == nil {
			adUser := u.GetUserADPropertysShort(recipient.Name)
			if adUser["findedInAD"] == false {
				adUser["displayName"] = recipient.Fio
				adUser["mail"] = recipient.Email
			}
			problem.Recipient = adUser
		}

	}

	ADiniciators := []map[string]interface{}{}
	ADexecutors := []map[string]interface{}{}
	ADobservers := []map[string]interface{}{}

	users, err := u.glpi.GetProblemUsers(id)
	if err == nil {
		for _, user := range users {

			adUser := u.GetUserADPropertysShort(user.Name)
			if adUser["findedInAD"] == false {
				adUser["displayName"] = user.Fio
				adUser["mail"] = user.Email
			}
			switch user.Type {
			case 1:
				if userRequesterInGLPI.Id == user.Id {
					problem.UserPresence.Initiator = true
				}
				ADiniciators = append(ADiniciators, adUser)
			case 2:
				if userRequesterInGLPI.Id == user.Id {
					problem.UserPresence.Executor = true
				}
				ADexecutors = append(ADexecutors, adUser)
			case 3:
				if userRequesterInGLPI.Id == user.Id {
					problem.UserPresence.Observer = true
				}
				ADobservers = append(ADobservers, adUser)
			}
		}
	}

	problem.Users.Initiators = ADiniciators
	problem.Users.Executors = ADexecutors
	problem.Users.Observers = ADobservers

	groups, err := u.glpi.GetProblemGroupForCurrentUser(id, userRequesterInGLPI.Id)

	Giniciators := []entity.GLPIGroup{}
	Gexetutors := []entity.GLPIGroup{}
	Gobservers := []entity.GLPIGroup{}
	if err == nil {
		for _, group := range groups {
			switch group.Type {
			case 1:
				if group.Presence {
					problem.GroupPresence.Initiator = true

				}
				Giniciators = append(Giniciators, group)
			case 2:
				if group.Presence {
					problem.GroupPresence.Executor = true

				}
				Gexetutors = append(Gexetutors, group)
			case 3:
				if group.Presence {
					problem.GroupPresence.Observer = true

				}
				Gobservers = append(Gobservers, group)
			}
		}
		problem.Groups.Initiators = Giniciators
		problem.Groups.Executors = Gexetutors
		problem.Groups.Observers = Gobservers

	}

	tickets, err := u.glpi.GetProblemTickets(id)
	problem.Krit = 0
	if err == nil {
		problem.Tickets = tickets
		for _, ticket := range tickets {
			if ticket.Krit > problem.Krit {
				problem.Krit = ticket.Krit
			}
			if ticket.Category > problem.Category {
				problem.Category = ticket.Category
			}
		}
	}

	return problem, nil
}
