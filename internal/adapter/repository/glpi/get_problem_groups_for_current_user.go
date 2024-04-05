package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetProblemGroupForCurrentUser(id string, user int) (users []entity.GLPIGroup, err error) {

	sql := fmt.Sprintf(`
            select glpi_groups_problems.groups_id AS 'id', glpi_groups.name, glpi_groups_problems.type,
			case when (SELECT users_id FROM glpi_groups_users WHERE groups_id=glpi_groups_problems.groups_id AND users_id=%d) IS NULL then 'false'
    	    ELSE 'true'
 		    end as presence
			FROM glpi_groups_problems
			INNER JOIN glpi_groups ON glpi_groups.id=glpi_groups_problems.groups_id
			WHERE glpi_groups_problems.problems_id = %s`, user, id)
	_, err = r.db.Select(&users, sql)
	return users, err

}
