package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetTicketWorks(ticketID string) (work []entity.GLPI_Work, err error) {
	var proc = fmt.Sprintf(`
	SELECT CONCAT('c-',glpi_itilfollowups.id) AS id , glpi_itilfollowups.content, is_private, glpi_itilfollowups.date_creation, glpi_itilfollowups.date_mod, name, CONCAT(realname," ", firstname) AS author, "commens" AS type
	FROM glpi_itilfollowups
	LEFT JOIN glpi_users ON glpi_itilfollowups.users_id= glpi_users.id
	WHERE items_id=%[1]s AND itemtype="Ticket"
	UNION
	SELECT CONCAT('r-',glpi_itilsolutions.id) AS id, glpi_itilsolutions.content, 0 as is_private, glpi_itilsolutions.date_creation, glpi_itilsolutions.date_mod, name, CONCAT(realname," ", firstname) AS author, "solutions" AS type
	FROM glpi_itilsolutions
	LEFT JOIN glpi_users ON glpi_itilsolutions.users_id= glpi_users.id
	WHERE items_id=%[1]s AND itemtype="Ticket"
	UNION
	SELECT CONCAT('t-',glpi_tickettasks.id) AS id, glpi_tickettasks.content, 0 as is_private, glpi_tickettasks.date_creation, glpi_tickettasks.date_mod, name, CONCAT(realname," ", firstname) AS author, "tasks" AS type
	FROM glpi_tickettasks
	LEFT JOIN glpi_users ON glpi_tickettasks.users_id= glpi_users.id
	WHERE tickets_id=%[1]s
	UNION
	SELECT CONCAT('ti-',glpi_tickets.id) AS id, glpi_tickets.content, 0 as is_private, glpi_tickets.date_creation, glpi_tickets.date_mod,"-" AS NAME,
	(SELECT user_name FROM glpi_logs WHERE itemtype="Ticket" and items_id=%[1]s order by id desc LIMIT 1) AS author, "create" AS type
	  from glpi_tickets WHERE id=%[1]s`, ticketID)
	_, err = r.db.Select(&work, proc)
	if err != nil {
		return nil, err
	}
	return work, nil
}
