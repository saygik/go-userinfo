package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

// ************************* Все отказы **********************************//
func (r *Repository) GetOtkazes(startdate string, enddate string) (otkazes []entity.GLPI_Otkaz, err error) {

	sql := fmt.Sprintf(
		`SELECT glpi_tickets.id as 'id',
			CASE
				WHEN STATUS<5 AND glpi_plugin_fields_failcategoryfielddropdowns.id>8
					THEN 3
				WHEN STATUS>=5 AND glpi_plugin_fields_failcategoryfielddropdowns.id>8
					THEN 2
				WHEN STATUS<5 AND (glpi_plugin_fields_failcategoryfielddropdowns.id<9 or ISNULL(glpi_plugin_fields_failcategoryfielddropdowns.id))
					THEN 1
				ELSE 0
			END AS 'krit',
			glpi_plugin_fields_failcategoryfielddropdowns.id as 'category',
			glpi_tickets.name AS 'name',
            glpi_tickets.is_deleted AS 'deleted',
			glpi_tickets.status
				AS 'status',
			glpi_tickets.impact
				AS 'impact',
			glpi_entities.completename as 'company',

			glpi_tickets.date AS 'date',
			IFNULL(glpi_tickets.solvedate, '') AS 'solvedate',
			IFNULL((SELECT CONCAT('[',GROUP_CONCAT(CONCAT('{"id":', glpi_problems_tickets.id, ', "name":',glpi_problems_tickets.problems_id,'}') ),']') as countpr FROM glpi_problems_tickets WHERE glpi_problems_tickets.tickets_id= glpi_tickets.id GROUP BY glpi_problems_tickets.tickets_id),'[]') AS 'problems',
			glpi_tickets.content AS 'content'
			FROM glpi_tickets
			LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			LEFT JOIN glpi_plugin_fields_ticketfailures ON glpi_plugin_fields_ticketfailures.items_id=glpi_tickets.id
			LEFT JOIN glpi_plugin_fields_failcategoryfielddropdowns ON glpi_plugin_fields_failcategoryfielddropdowns.id=glpi_plugin_fields_ticketfailures.plugin_fields_failcategoryfielddropdowns_id
			WHERE glpi_tickets.is_deleted<>TRUE  AND glpi_plugin_fields_failcategoryfielddropdowns.id>4
			AND
			((date>='%[1]s' AND date <='%[2]s') OR (solvedate>='%[1]s' AND solvedate <='%[2]s') OR (date<'%[1]s' AND solvedate >'%[2]s') OR (date<'%[1]s' AND solvedate is null))
			ORDER BY date desc
		`, startdate, enddate)
	_, err = r.db.Select(&otkazes, sql)
	return otkazes, err
}
