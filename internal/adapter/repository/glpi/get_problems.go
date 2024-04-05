package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

// ************************* Все отказы **********************************//
func (r *Repository) GetProblems(startdate string, enddate string) (problems []entity.GLPI_Problem, err error) {

	sql := fmt.Sprintf(
		`SELECT glpi_problems.id AS 'id',
        CASE  WHEN STATUS<5 AND glpi_plugin_fields_failcategoryfielddropdowns.id>7 THEN 3
               WHEN STATUS>=5 AND glpi_plugin_fields_failcategoryfielddropdowns.id>7 THEN 2
               WHEN STATUS<5 AND (glpi_plugin_fields_failcategoryfielddropdowns.id<8 or ISNULL(glpi_plugin_fields_failcategoryfielddropdowns.id)) THEN 1
               ELSE 0
        END AS 'krit',
        IFNULL(glpi_plugin_fields_failcategoryfielddropdowns.id,0) AS 'category',
        glpi_problems.name as 'name', glpi_problems.content AS 'content', status as 'status', impact AS 'impact', DATE AS 'date' , IFNULL(glpi_problems.solvedate, '') AS 'solvedate',
        CONCAT (TIMESTAMPDIFF(HOUR, DATE, IFNULL(solvedate, NOW())), "ч. ", TIMESTAMPDIFF(MINUTE, DATE, IFNULL(solvedate, NOW()) - INTERVAL TIMESTAMPDIFF(HOUR, DATE, IFNULL(solvedate, NOW())) HOUR )," мин.") AS 'solvetime',
        glpi_entities.completename  as 'company',
        IFNULL(pt.ticket_count,0) AS 'ticketscount',
		IFNULL((SELECT CONCAT('[',GROUP_CONCAT(CONCAT('{"id":', glpi_problems_tickets.id, ', "name":',glpi_problems_tickets.tickets_id,'}') ),']') as countpr FROM glpi_problems_tickets WHERE glpi_problems_tickets.problems_id= glpi_problems.id GROUP BY glpi_problems_tickets.problems_id),'[]') AS 'ticketsid'
        FROM glpi_problems
        LEFT JOIN glpi_entities ON glpi_problems.entities_id = glpi_entities.id
        LEFT JOIN glpi_itilcategories ON glpi_problems.itilcategories_id = glpi_itilcategories.id
        LEFT JOIN (SELECT problems_id, count(id) AS ticket_count FROM glpi_problems_tickets GROUP BY problems_id) pt ON pt.problems_id=glpi_problems.id
        LEFT JOIN
		(SELECT glpi_problems_tickets.problems_id, IFNULL(MAX(glpi_plugin_fields_ticketfailures.plugin_fields_failcategoryfielddropdowns_id),1) AS FailCategory
		FROM  glpi_problems_tickets
		LEFT JOIN glpi_plugin_fields_ticketfailures ON glpi_plugin_fields_ticketfailures.items_id=glpi_problems_tickets.tickets_id
		GROUP BY problems_id) pcat ON glpi_problems.id=pcat.problems_id
		LEFT JOIN glpi_plugin_fields_failcategoryfielddropdowns ON glpi_plugin_fields_failcategoryfielddropdowns.id=pcat.FailCategory
		WHERE is_deleted=0 and glpi_problems.name not like '%%2222%%' and glpi_problems.name not like '%%2222%%'
		AND
			((date>='%[1]s' AND date <='%[2]s') OR (solvedate>='%[1]s' AND solvedate <='%[2]s') OR (date<'%[1]s' AND solvedate >'%[2]s') OR (date<'%[1]s' AND solvedate is null))
			ORDER BY date desc
		`, startdate, enddate)
	_, err = r.db.Select(&problems, sql)
	return problems, err
}
