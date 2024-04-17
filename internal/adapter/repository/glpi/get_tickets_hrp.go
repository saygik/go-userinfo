package glpi

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetHRPTickets() (tickets []entity.GLPI_Ticket, err error) {

	sql := `SELECT glpi_tickets.id AS 'id', glpi_tickets.name AS 'name', content  AS 'content', glpi_entities.completename as 'company' FROM glpi_tickets
	        INNER JOIN glpi_entities ON glpi_tickets.entities_id=glpi_entities.id
            WHERE requesttypes_id=8 AND status<5  AND glpi_tickets.id not IN (select tickets_id AS id FROM glpi_tickets_hrp) AND glpi_entities.completename LIKE 'БЖД > ИВЦ2%'`
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}
