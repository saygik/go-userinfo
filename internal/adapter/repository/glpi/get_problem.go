package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetProblem(id string) (problem entity.GLPI_Problem, err error) {

	sql := fmt.Sprintf(`
				SELECT glpi_problems.id,
				glpi_problems.name,
				glpi_problems.status,
				glpi_problems.is_deleted,
				glpi_problems.is_recursive as 'recursive',
				IFNULL(glpi_problems.users_id_recipient,0) as 'recipient_id',
				IFNULL(glpi_problems.content, '') as 'content',
				glpi_problems.impact,
				IFNULL(glpi_problems.impactcontent, '') as 'impactcontent',
				IFNULL(glpi_problems.causecontent, '') as 'causecontent',
				IFNULL(glpi_problems.symptomcontent, '') as 'symptomcontent',
				glpi_entities.completename AS 'company',
				glpi_problems.date AS 'date' ,
				IFNULL(glpi_problems.solvedate,'') as 'solvedate',
				glpi_problems.date_mod AS 'datemod',
				IFNULL(glpi_problems.closedate,'') as 'closedate',
				glpi_problems.date_creation AS 'date_creation',
				(SELECT count(id) FROM glpi_problems_tickets where glpi_problems_tickets.problems_id=glpi_problems.id)   AS 'ticketscount'
				FROM glpi_problems
				LEFT JOIN glpi_entities ON glpi_problems.entities_id = glpi_entities.id
				WHERE glpi_problems.id=%s`, id)
	err = r.db.SelectOne(&problem, sql)

	return problem, err

}
