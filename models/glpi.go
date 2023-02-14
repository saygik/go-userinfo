package models

import (
	"fmt"

	"github.com/saygik/go-userinfo/glpidb"
)

type GLPIUser struct {
	Id       int64             `db:"id" json:"id"`
	Name     string            `db:"name" json:"name"`
	Self     string            `db:"self" json:"self"`
	Date     string            `db:"date" json:"date"`
	Author   int64             `db:"author" json:"author"`
	Executor int64             `db:"executor" json:"executor"`
	Profiles []GLPIUserProfile `db:"profiles" json:"profiles"`
}

type GLPIUserProfile struct {
	Id        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	EName     string `db:"ename" json:"ename"`
	Recursive bool   `db:"recursive" json:"recursive"`
}

type Software struct {
	Id             int64  `db:"id" json:"id"`
	Name           string `db:"name" json:"name"`
	Ename          string `db:"ename" json:"company"`
	Comment        string `db:"comment" json:"comment"`
	Locations      string `db:"locations" json:"locations,omitempty"`
	Manufacture    string `db:"manufacture" json:"manufacture"`
	Description1   string `db:"description1" json:"description1"`
	Description2   string `db:"description2" json:"description2"`
	Murl           string `db:"murl" json:"manual_url"`
	Durl           string `db:"durl" json:"doc_url"`
	IsRecursive    int64  `db:"is_recursive" json:"is_recursive"`
	Groups_id_tech int64  `db:"groups_id_tech" json:"groups_id_tech"`
	Users_id_tech  int64  `db:"users_id_tech" json:"users_id_tech"`
	Extauth        int64  `db:"extauth" json:"ext_auth"`
	Clients        int64  `db:"clients" json:"clients"`
}

type GLPIModel struct{}

func (m GLPIModel) GetUserByName(login string) (user GLPIUser, err error) {
	sql := fmt.Sprintf(
		`SELECT u.id , u.name, IFNULL(u.last_login,'-') AS date, IFNULL((SELECT glpi_entities.completename  FROM glpi_profiles_users 
                INNER JOIN glpi_entities ON glpi_entities.id=glpi_profiles_users.entities_id 
                WHERE glpi_profiles_users.users_id=u.id AND glpi_profiles_users.profiles_id=1 LIMIT 1),'-') AS self,
                (SELECT COUNT(glpi_tickets_users.id) FROM glpi_tickets_users INNER JOIN glpi_tickets ON glpi_tickets.id=glpi_tickets_users.tickets_id WHERE glpi_tickets_users.users_id=u.id AND glpi_tickets_users.type=1 AND is_deleted=0) AS author,
                (SELECT COUNT(glpi_tickets_users.id) FROM glpi_tickets_users INNER JOIN glpi_tickets ON glpi_tickets.id=glpi_tickets_users.tickets_id WHERE glpi_tickets_users.users_id=u.id AND glpi_tickets_users.type=2 AND is_deleted=0) AS executor 
                 FROM (SELECT * FROM glpi_users  WHERE glpi_users.name= '%s' ) u`, login)
	err = glpidb.GetDB().SelectOne(&user, sql)
	return user, err
}

func (m GLPIModel) GetUserProfiles(id int64) (profiles []GLPIUserProfile, err error) {
	sql := fmt.Sprintf(
		`SELECT glpi_profiles_users.profiles_id AS id, glpi_profiles.name AS 'name', glpi_entities.completename AS ename, glpi_profiles_users.is_recursive AS 'recursive' 
		FROM  glpi_profiles_users
		INNER JOIN glpi_profiles ON glpi_profiles_users.profiles_id=glpi_profiles.id
		INNER JOIN glpi_entities ON glpi_profiles_users.entities_id=glpi_entities.id
		WHERE users_id=%d`, id)
	_, err = glpidb.GetDB().Select(&profiles, sql)
	return profiles, err
}
func (m GLPIModel) GetSoftwares() (softwares []Software, err error) {
	sql := fmt.Sprintf(
		`SELECT glpi_softwares.id, glpi_softwares.name, glpi_softwares.comment, IFNULL(glpi_entities.completename,'') AS ename, glpi_softwares.is_recursive, 
		IFNULL(glpi_locations.completename,'') AS locations, glpi_softwares.groups_id_tech, glpi_softwares.users_id_tech, IFNULL(glpi_manufacturers.name,'') AS manufacture,
		IFNULL(softadd.descriptionfieldtwo,'') AS description1, IFNULL(softadd.moredescriptionfield,'') AS description2, IFNULL(softadd.externalauthenticationfieldtwo,0) AS extauth, 
		IFNULL(softadd.clientsoftwarefieldtwo,0) AS clients, IFNULL(softadd.servicemanualurlfieldtwo,'') murl, IFNULL(softadd.technicaldescriptionurlfield,'') AS durl 
		from glpi_softwares
		INNER JOIN glpi_entities ON glpi_softwares.entities_id=glpi_entities.id
		LEFT JOIN glpi_locations ON glpi_softwares.locations_id=glpi_locations.id
		LEFT JOIN glpi_manufacturers ON glpi_softwares.manufacturers_id=glpi_manufacturers.id
		LEFT JOIN glpi_plugin_fields_softwareadditionals softadd ON glpi_softwares.id=softadd.items_id
		WHERE glpi_softwares.is_deleted=0`)
	_, err = glpidb.GetDB().Select(&softwares, sql)
	return softwares, err
}
