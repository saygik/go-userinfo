package usecase

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetGLPITicketSolutionTemplates(ticketID string) ([]entity.GLPI_Ticket_Profile, error) {
	profiles, err := u.glpi.GetGLPITicketSolutionTemplates(ticketID)
	if err != nil {
		return profiles, u.Error(fmt.Sprintf("ошибка MySQL: %s", err.Error()))
	}
	return profiles, nil
}
