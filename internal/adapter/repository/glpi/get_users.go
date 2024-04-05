package glpi

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (r *Repository) GetUsers() (users []entity.GLPIUserShort, err error) {
	sql := `SELECT  u.id as id, u.name as name, IFNULL(IF(realname='' AND firstname='', u.name , CONCAT(realname, ' ',firstname)),'-') AS realname,
		authtype, LOWER(IFNULL(glpi_authldaps.name,'-')) AS ad, IFNULL(last_login, '-') as last_login,
		(SELECT COUNT(glpi_tickets_users.id) FROM glpi_tickets_users INNER JOIN glpi_tickets ON glpi_tickets.id=glpi_tickets_users.tickets_id WHERE glpi_tickets_users.users_id=u.id AND glpi_tickets_users.type=1 AND is_deleted=0) AS author
		FROM glpi_users u LEFT JOIN glpi_authldaps ON glpi_authldaps.id=u.auths_id WHERE is_deleted=false`
	_, err = r.db.Select(&users, sql)
	return users, err
}
