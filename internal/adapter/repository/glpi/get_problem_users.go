package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetProblemUsers(id string) (users []entity.GLPIUser, err error) {

	sql := fmt.Sprintf(`
          SELECT glpi_users.id AS 'id', name, CONCAT(realname," ", firstname) AS 'fio', IFNULL((SELECT email from glpi_useremails WHERE users_id=glpi_users.id LIMIT 1),'') AS 'email', tu.type AS 'type'
          FROM (SELECT users_id,type FROM glpi_problems_users WHERE problems_id=%s) tu
          INNER JOIN glpi_users on tu.users_id=glpi_users.id`, id)
	_, err = r.db.Select(&users, sql)
	return users, err

}
