package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

// ************************* Незакрытые заявки **********************************//

func (r *Repository) GetTicketsNonClosedFromIniciator(id int) (tickets []entity.GLPI_Ticket, err error) {
	sql := fmt.Sprintf(`
		SELECT IFNULL(JSON_EXTRACT(ancestors_cache, '$.*'),JSON_ARRAY(0)) AS 'orgs',IFNULL(fkat,0) as category,
		gt.id , gt.content, gt.status, gt.name, gt.impact, glpi_entities.completename as company, glpi_entities.id as eid, IFNULL(gt.date,'') as date, gt.date_mod, gt.date_creation, IFNULL(gt.solvedate,'') as solvedate,
		CONCAT(ifnull(NULLIF(glpi_users.realname, ''), 'не опреденен'),' ', ifnull(NULLIF(glpi_users.firstname, ''),'')) AS author,
		IFNULL((SELECT CONCAT('[',GROUP_CONCAT(CONCAT('{"id":', glpi_tickets_users.users_id, ', "name":"',glpi_users.name,'"', ', "type":',glpi_tickets_users.type,'}') ),']')  FROM glpi_tickets_users INNER JOIN glpi_users ON glpi_users.id=glpi_tickets_users.users_id  WHERE glpi_tickets_users.tickets_id = gt.id  ),'[]') AS users_s,
		IFNULL((SELECT CONCAT('[',GROUP_CONCAT(CONCAT('{"id":', glpi_groups_tickets.groups_id, ', "name":"',glpi_groups.name,'"', ', "type":',glpi_groups_tickets.type,'}') ),']')  FROM glpi_groups_tickets INNER JOIN glpi_groups ON glpi_groups.id=glpi_groups_tickets.groups_id  WHERE glpi_groups_tickets.tickets_id = gt.id  ),'[]') AS user_groups_s
		FROM (SELECT t.* FROM (SELECT *  from glpi_tickets WHERE glpi_tickets.is_deleted=0  AND STATUS <5) t
        INNER JOIN glpi_tickets_users ON t.id=glpi_tickets_users.tickets_id
        WHERE glpi_tickets_users.type=1 AND users_id=%d) gt
		INNER JOIN glpi_entities ON gt.entities_id=glpi_entities.id
		LEFT JOIN glpi_users ON gt.users_id_recipient=glpi_users.id
		LEFT JOIN  (SELECT items_id,plugin_fields_failcategoryfielddropdowns_id AS fkat  from glpi_plugin_fields_ticketfailures WHERE plugin_fields_failcategoryfielddropdowns_id>4) fc ON fc.items_id=gt.id
		`, id)
	_, err = r.db.Select(&tickets, sql)
	return tickets, err
}
