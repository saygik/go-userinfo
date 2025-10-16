package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetTicket(id string) (ticket entity.GLPI_Ticket, err error) {

	sql := fmt.Sprintf(`SELECT
        IFNULL(JSON_EXTRACT(ancestors_cache, '$.*'),JSON_ARRAY(0)) AS 'orgs',
        IFNULL(fkat,0) as category,
			CASE
				WHEN STATUS<5 AND fkat>8
					THEN 3
				WHEN STATUS>=5 AND fkat>8
					THEN 2
				WHEN STATUS<5 AND (fkat<9 or ISNULL(fkat))
					THEN 1
				ELSE 0
			END AS 'krit',
  		gt.id ,
	    gt.content,
	    gt.status,
	    gt.name,
		gt.impact,
	    glpi_entities.completename as company,
        glpi_entities.id as eid,
		IFNULL(gt.date,'') AS 'date',
		gt.date_mod AS 'date_mod',
		IFNULL(gt.closedate,'') AS 'closedate',
		gt.date_creation AS 'date_creation',
		IFNULL(gt.solvedate,'') AS 'solvedate',
		gt.is_deleted,
		IFNULL((SELECT users_id FROM glpi_tickets_users WHERE tickets_id=%[1]s AND  glpi_tickets_users.type=1 LIMIT 1),0) as 'recipient_id',
		gt.type,
		gt.requesttypes_id,
		IFNULL((SELECT CONCAT('[',GROUP_CONCAT(CONCAT('{"id":', glpi_tickets_users.users_id, ', "name":"',glpi_users.name,'"', ', "type":',glpi_tickets_users.type,'}') ),']')  FROM glpi_tickets_users INNER JOIN glpi_users ON glpi_users.id=glpi_tickets_users.users_id  WHERE glpi_tickets_users.tickets_id = gt.id  ),'[]') AS users_s,
		IFNULL((SELECT CONCAT('[',GROUP_CONCAT(CONCAT('{"id":', glpi_groups_tickets.groups_id, ', "name":"',glpi_groups.name,'"', ', "type":',glpi_groups_tickets.type,'}') ),']')  FROM glpi_groups_tickets INNER JOIN glpi_groups ON glpi_groups.id=glpi_groups_tickets.groups_id  WHERE glpi_groups_tickets.tickets_id = gt.id  ),'[]') AS user_groups_s,
        IFNULL((SELECT IFNULL(completename,"")  FROM glpi_groups_tickets INNER JOIN glpi_groups ON glpi_groups.id=glpi_groups_tickets.groups_id  WHERE type=2 and glpi_groups_tickets.tickets_id = gt.id),'') AS 'group_name',
     	(SELECT count(id) FROM glpi_problems_tickets where glpi_problems_tickets.tickets_id=gt.id)   AS 'problemscount'
		FROM (SELECT * FROM glpi_tickets WHERE glpi_tickets.id=%[1]s) gt
		INNER JOIN glpi_entities ON gt.entities_id=glpi_entities.id
		LEFT JOIN  (SELECT items_id,plugin_fields_failcategoryfielddropdowns_id AS fkat  from glpi_plugin_fields_ticketfailures) fc ON fc.items_id=gt.id
	`, id)
	err = r.db.SelectOne(&ticket, sql)

	return ticket, err

}
