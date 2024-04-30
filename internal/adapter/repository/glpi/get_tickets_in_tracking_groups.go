package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

// ************************* Незакрытые заявки **********************************//

func (r *Repository) GetTicketsInExecutionGroups(ids string) (tickets []entity.GLPI_Ticket, err error) {
	sql := fmt.Sprintf(`SELECT IFNULL(fkat,0) as category,		gt.id , gt.content, gt.status, gt.name, gt.impact, glpi_entities.completename as company, IFNULL(gt.date,'') as DATE,tg.id AS group_id,tg.name AS group_name
		FROM (SELECT * FROM glpi_tickets WHERE glpi_tickets.is_deleted=0  AND STATUS <5 ) gt
		INNER JOIN glpi_entities ON gt.entities_id=glpi_entities.id
		LEFT JOIN  (SELECT items_id,plugin_fields_failcategoryfielddropdowns_id AS fkat  from glpi_plugin_fields_ticketfailures WHERE plugin_fields_failcategoryfielddropdowns_id>4) fc ON fc.items_id=gt.id
		LEFT JOIN (SELECT glpi_groups_tickets.groups_id as 'id', glpi_groups.name, glpi_groups_tickets.tickets_id
		FROM glpi_groups_tickets
		INNER JOIN glpi_groups ON glpi_groups.id=glpi_groups_tickets.groups_id
		WHERE glpi_groups_tickets.type=2) tg ON gt.id=tg.tickets_id
        WHERE tg.id IN (%s)`, ids)
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}
