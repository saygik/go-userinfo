package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

// ************************* Все отказы **********************************//
func (r *Repository) GetProblemTickets(id string) (otkazes []entity.GLPI_Otkaz, err error) {

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
			IFNULL(glpi_plugin_fields_failcategoryfielddropdowns.id,0) as 'category',
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
			WHERE glpi_tickets.id IN (SELECT tickets_id AS 'id' FROM glpi_problems_tickets where glpi_problems_tickets.problems_id=%s)`, id)
	_, err = r.db.Select(&otkazes, sql)
	return otkazes, err
}
