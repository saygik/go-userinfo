package models

import (
	"fmt"
	"github.com/saygik/go-userinfo/glpidb"
)

type GLPIUser struct {
	Name     string `db:"name" json:"name"`
	Self     string `db:"self" json:"self"`
	Date     string `db:"date" json:"date"`
	Author   int64  `db:"author" json:"author"`
	Executor int64  `db:"executor" json:"executor"`
}

type GLPIModel struct{}

func (m GLPIModel) GetUserByName(login string) (user GLPIUser, err error) {
	sql := fmt.Sprintf(
		`SELECT u.name, IFNULL(u.last_login,'-') AS date, IFNULL((SELECT glpi_entities.completename  FROM glpi_profiles_users 
                INNER JOIN glpi_entities ON glpi_entities.id=glpi_profiles_users.entities_id 
                WHERE glpi_profiles_users.users_id=u.id AND glpi_profiles_users.profiles_id=1 LIMIT 1),'-') AS self,
                (SELECT COUNT(glpi_tickets_users.id) FROM glpi_tickets_users INNER JOIN glpi_tickets ON glpi_tickets.id=glpi_tickets_users.tickets_id WHERE glpi_tickets_users.users_id=u.id AND glpi_tickets_users.type=1 AND is_deleted=0) AS author,
                (SELECT COUNT(glpi_tickets_users.id) FROM glpi_tickets_users INNER JOIN glpi_tickets ON glpi_tickets.id=glpi_tickets_users.tickets_id WHERE glpi_tickets_users.users_id=u.id AND glpi_tickets_users.type=2 AND is_deleted=0) AS executor 
                 FROM (SELECT * FROM glpi_users  WHERE glpi_users.name= '%s' ) u`, login)
	err = glpidb.GetDB().SelectOne(&user, sql)
	return user, err
}
