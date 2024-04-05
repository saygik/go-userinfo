package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetTicketGroupForCurrentUser(id string, user int) (users []entity.GLPIGroup, err error) {

	sql := fmt.Sprintf(`
            select glpi_groups_tickets.groups_id as 'id', glpi_groups.name, glpi_groups_tickets.type,
            case when (SELECT users_id FROM glpi_groups_users WHERE groups_id=glpi_groups_tickets.groups_id AND users_id=%d) IS NULL then 'false'
      		ELSE 'true'
     		end as 'presence'
			FROM glpi_groups_tickets
			INNER JOIN glpi_groups ON glpi_groups.id=glpi_groups_tickets.groups_id
			WHERE glpi_groups_tickets.tickets_id = %s`, user, id)
	_, err = r.db.Select(&users, sql)
	return users, err

}
