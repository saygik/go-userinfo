package glpi

import (
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetUserByName(login string) (user entity.GLPIUser, err error) {
	sql := fmt.Sprintf(
		`SELECT u.id , u.name, IFNULL(u.last_login,'-') AS date, e.completename as org,  not isnull(api_token) AS api, IFNULL((SELECT glpi_entities.completename  FROM glpi_profiles_users
                INNER JOIN glpi_entities ON glpi_entities.id=glpi_profiles_users.entities_id
                WHERE glpi_profiles_users.users_id=u.id AND glpi_profiles_users.profiles_id=1 LIMIT 1),'-') AS self,
                IFNULL((SELECT glpi_entities.id  FROM glpi_profiles_users
				INNER JOIN glpi_entities ON glpi_entities.id=glpi_profiles_users.entities_id
                WHERE glpi_profiles_users.users_id=u.id AND glpi_profiles_users.profiles_id=1 LIMIT 1),0) AS selfid,
                (SELECT COUNT(glpi_tickets_users.id) FROM glpi_tickets_users INNER JOIN glpi_tickets ON glpi_tickets.id=glpi_tickets_users.tickets_id WHERE glpi_tickets_users.users_id=u.id AND glpi_tickets_users.type=1 AND is_deleted=0) AS author,
                (SELECT COUNT(glpi_tickets_users.id) FROM glpi_tickets_users INNER JOIN glpi_tickets ON glpi_tickets.id=glpi_tickets_users.tickets_id WHERE glpi_tickets_users.users_id=u.id AND glpi_tickets_users.type=2 AND is_deleted=0) AS executor
                 FROM (SELECT * FROM glpi_users  WHERE glpi_users.name= '%s' ) u INNER JOIN glpi_entities e ON e.id=u.entities_id`, login)
	err = r.db.SelectOne(&user, sql)
	return user, err
}
