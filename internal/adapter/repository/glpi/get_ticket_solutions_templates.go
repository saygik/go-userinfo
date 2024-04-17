package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetGLPITicketSolutionTemplates(ticketID string) (profiles []entity.GLPI_Ticket_Profile, err error) {
	var proc = fmt.Sprintf(`
		SELECT glpi_solutiontemplates.id,glpi_solutiontemplates.name,glpi_solutiontemplates.content,glpi_entities.completename as 'company'
		from glpi_solutiontemplates
		LEFT JOIN glpi_entities ON glpi_solutiontemplates.entities_id = glpi_entities.id
		WHERE (INSTR((
   		SELECT	glpi_entities.completename as 'company'
		FROM glpi_tickets
   		LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
   		WHERE glpi_tickets.id=%[1]s),glpi_entities.completename)>0 AND glpi_solutiontemplates.is_recursive=1) or
   		(SELECT	glpi_entities.completename as 'company'
		FROM glpi_tickets
   		LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
   		WHERE glpi_tickets.id=%[1]s)=glpi_entities.completename
   		ORDER BY name
	`, ticketID)
	_, err = r.db.Select(&profiles, proc)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}
