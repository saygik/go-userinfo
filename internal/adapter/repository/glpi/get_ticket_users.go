package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetTicketUsers(id string) (users []entity.GLPIUser, err error) {

	sql := fmt.Sprintf(`
           SELECT glpi_users.id, name, CONCAT(realname," ", firstname) AS fio, IFNULL((SELECT email from glpi_useremails WHERE users_id=glpi_users.id LIMIT 1),'') as email, tu.type AS 'type'
         FROM (SELECT * from glpi_tickets_users  WHERE tickets_id=%s) tu
			INNER JOIN glpi_users on tu.users_id=glpi_users.id`, id)
	_, err = r.db.Select(&users, sql)
	return users, err

}
