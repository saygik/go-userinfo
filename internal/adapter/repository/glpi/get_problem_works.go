package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetProblemWorks(id string) (work []entity.GLPI_Work, err error) {
	var proc = fmt.Sprintf(`
SELECT CONCAT('c-',glpi_itilfollowups.id) AS id , glpi_itilfollowups.content, is_private, glpi_itilfollowups.date_creation, glpi_itilfollowups.date_mod, name, CONCAT(realname," ", firstname) AS author, "commens" AS type, 0 as "status"
	FROM glpi_itilfollowups
	LEFT JOIN glpi_users ON glpi_itilfollowups.users_id= glpi_users.id
	WHERE items_id=%[1]s AND itemtype="Problem"
	UNION
	SELECT CONCAT('r-',glpi_itilsolutions.id) AS id, glpi_itilsolutions.content, 0 as is_private, glpi_itilsolutions.date_creation, glpi_itilsolutions.date_mod, name, CONCAT(realname," ", firstname) AS author, "solutions" AS type, status
	FROM glpi_itilsolutions
	LEFT JOIN glpi_users ON glpi_itilsolutions.users_id= glpi_users.id
	WHERE items_id=%[1]s AND itemtype="Problem"
	UNION
	SELECT CONCAT('t-',glpi_problemtasks.id) AS id, glpi_problemtasks.content, 0 as is_private, glpi_problemtasks.date_creation, glpi_problemtasks.date_mod, name, CONCAT(realname," ", firstname) AS author, "tasks" AS type, state as "status"
	FROM glpi_problemtasks
	LEFT JOIN glpi_users ON glpi_problemtasks.users_id= glpi_users.id
	WHERE problems_id=%[1]s
	UNION
	SELECT CONCAT('ti-',glpi_problems.id) AS id, glpi_problems.content, 0 as is_private, glpi_problems.date_creation, glpi_problems.date_mod,"-" AS NAME,
	(SELECT user_name FROM glpi_logs WHERE itemtype="Problem" and items_id=%[1]s order by id desc LIMIT 1) AS author, "create" AS type, 0 as "status"
	  from glpi_problems WHERE id=%[1]s`, id)
	_, err = r.db.Select(&work, proc)
	if err != nil {
		return nil, err
	}
	return work, nil
}
