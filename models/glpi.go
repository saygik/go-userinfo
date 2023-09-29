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
	Org      string            `db:"org" json:"org"`
	Author   int64             `db:"author" json:"author"`
	Executor int64             `db:"executor" json:"executor"`
	Profiles []GLPIUserProfile `db:"profiles" json:"profiles"`
}
type GLPIUserShort struct {
	Id        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Realname  string `db:"realname" json:"realname"`
	Authtype  int64  `db:"authtype" json:"authtype"`
	Deleted   int64  `db:"is_deleted" json:"is_deleted"`
	AD        string `db:"ad" json:"ad"`
	LastLogin string `db:"last_login" json:"last_login"`
	Author    int64  `db:"author" json:"author"`
}
type GLPIUserProfile struct {
	Id        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	EName     string `db:"ename" json:"ename"`
	Eid       int64  `db:"eid" json:"eid"`
	Orgs      string `db:"orgs" json:"orgs"`
	Recursive bool   `db:"recursive" json:"recursive"`
}
type SoftwareAdmins struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
type Software struct {
	Id             int64    `db:"id" json:"id"`
	Name           string   `db:"name" json:"name"`
	Ename          string   `db:"ename" json:"company"`
	Comment        string   `db:"comment" json:"comment"`
	Locations      string   `db:"locations" json:"locations,omitempty"`
	Manufacture    string   `db:"manufacture" json:"manufacture"`
	Description1   string   `db:"description1" json:"description1"`
	Description2   string   `db:"description2" json:"description2"`
	Murl           string   `db:"murl" json:"manual_url"`
	Durl           string   `db:"durl" json:"icon_url"`
	IsRecursive    int64    `db:"is_recursive" json:"is_recursive"`
	Groups_id_tech int64    `db:"groups_id_tech" json:"groups_id_tech"`
	Users_id_tech  int64    `db:"users_id_tech" json:"users_id_tech"`
	Extauth        int64    `db:"extauth" json:"ext_auth"`
	Clients        int64    `db:"clients" json:"clients"`
	GroupName      string   `db:"group_name" json:"group_name"`
	Admins         []string `json:"tech_users"`
}
type TicketsStats struct {
	Сount int64 `db:"count" json:"count"`
	Type  int64 `db:"type" json:"type"`
	Year  int64 `db:"year" json:"year"`
	Month int64 `db:"month" json:"month"`
}

type RegionsStats struct {
	Сount int64  `db:"count" json:"count"`
	Org   string `db:"org" json:"org"`
	Proc  int64  `db:"proc" json:"proc"`
}
type GLPIModel struct{}

func (m GLPIModel) GetUserByName(login string) (user GLPIUser, err error) {
	sql := fmt.Sprintf(
		`SELECT u.id , u.name, IFNULL(u.last_login,'-') AS date, e.completename as org, IFNULL((SELECT glpi_entities.completename  FROM glpi_profiles_users
                INNER JOIN glpi_entities ON glpi_entities.id=glpi_profiles_users.entities_id
                WHERE glpi_profiles_users.users_id=u.id AND glpi_profiles_users.profiles_id=1 LIMIT 1),'-') AS self,
                (SELECT COUNT(glpi_tickets_users.id) FROM glpi_tickets_users INNER JOIN glpi_tickets ON glpi_tickets.id=glpi_tickets_users.tickets_id WHERE glpi_tickets_users.users_id=u.id AND glpi_tickets_users.type=1 AND is_deleted=0) AS author,
                (SELECT COUNT(glpi_tickets_users.id) FROM glpi_tickets_users INNER JOIN glpi_tickets ON glpi_tickets.id=glpi_tickets_users.tickets_id WHERE glpi_tickets_users.users_id=u.id AND glpi_tickets_users.type=2 AND is_deleted=0) AS executor
                 FROM (SELECT * FROM glpi_users  WHERE glpi_users.name= '%s' ) u INNER JOIN glpi_entities e ON e.id=u.entities_id`, login)
	err = glpidb.GetDB().SelectOne(&user, sql)
	return user, err
}

func (m GLPIModel) GetUserProfiles(id int64) (profiles []GLPIUserProfile, err error) {
	sql := fmt.Sprintf(
		`SELECT glpi_profiles_users.profiles_id AS id, glpi_profiles.name AS 'name', glpi_entities.completename AS ename, glpi_profiles_users.is_recursive AS 'recursive',
		glpi_entities.id AS 'eid', IFNULL(JSON_EXTRACT(ancestors_cache, '$.*'),JSON_ARRAY(0)) AS 'orgs'
		FROM  glpi_profiles_users
		INNER JOIN glpi_profiles ON glpi_profiles_users.profiles_id=glpi_profiles.id
		INNER JOIN glpi_entities ON glpi_profiles_users.entities_id=glpi_entities.id
		WHERE users_id=%d`, id)
	_, err = glpidb.GetDB().Select(&profiles, sql)
	return profiles, err
}
func (m GLPIModel) GetSoftware(id int64) (software Software, err error) {
	sql := fmt.Sprintf(
		`SELECT glpi_softwares.id, glpi_softwares.name, glpi_softwares.comment, IFNULL(glpi_entities.completename,'') AS ename, glpi_softwares.is_recursive,
		IFNULL(glpi_locations.completename,'') AS locations, glpi_softwares.groups_id_tech, glpi_softwares.users_id_tech, IFNULL(glpi_manufacturers.name,'') AS manufacture,
		IFNULL(softadd.descriptionfieldtwo,'') AS description1, IFNULL(softadd.moredescriptionfield,'') AS description2, IFNULL(softadd.externalauthenticationfieldtwo,0) AS extauth,
		IFNULL(softadd.clientsoftwarefieldtwo,0) AS clients, IFNULL(softadd.servicemanualurlfieldtwo,'') murl, IFNULL(softadd.technicaldescriptionurlfield,'') AS durl, IFNULL(glpi_groups.name,'') as group_name
		from glpi_softwares
		INNER JOIN glpi_entities ON glpi_softwares.entities_id=glpi_entities.id
		LEFT JOIN glpi_locations ON glpi_softwares.locations_id=glpi_locations.id
		LEFT JOIN glpi_manufacturers ON glpi_softwares.manufacturers_id=glpi_manufacturers.id
		LEFT JOIN glpi_plugin_fields_softwareadditionals softadd ON glpi_softwares.id=softadd.items_id
		LEFT JOIN glpi_groups ON glpi_softwares.groups_id_tech=glpi_groups.id
		WHERE glpi_softwares.id=%d`, id)
	err = glpidb.GetDB().SelectOne(&software, sql)

	return software, err
}
func (m GLPIModel) GetSoftwares() (softwares []Software, err error) {
	sql := fmt.Sprintf(
		`SELECT glpi_softwares.id, glpi_softwares.name, glpi_softwares.comment, IFNULL(glpi_entities.completename,'') AS ename, glpi_softwares.is_recursive,
		IFNULL(glpi_locations.completename,'') AS locations, glpi_softwares.groups_id_tech, glpi_softwares.users_id_tech, IFNULL(glpi_manufacturers.name,'') AS manufacture,
		IFNULL(softadd.descriptionfieldtwo,'') AS description1, IFNULL(softadd.moredescriptionfield,'') AS description2, IFNULL(softadd.externalauthenticationfieldtwo,0) AS extauth,
		IFNULL(softadd.clientsoftwarefieldtwo,0) AS clients, IFNULL(softadd.servicemanualurlfieldtwo,'') murl, IFNULL(softadd.technicaldescriptionurlfield,'') AS durl, IFNULL(glpi_groups.name,'') as group_name
		from glpi_softwares
		INNER JOIN glpi_entities ON glpi_softwares.entities_id=glpi_entities.id
		LEFT JOIN glpi_locations ON glpi_softwares.locations_id=glpi_locations.id
		LEFT JOIN glpi_manufacturers ON glpi_softwares.manufacturers_id=glpi_manufacturers.id
		LEFT JOIN glpi_plugin_fields_softwareadditionals softadd ON glpi_softwares.id=softadd.items_id
		LEFT JOIN glpi_groups ON glpi_softwares.groups_id_tech=glpi_groups.id
		WHERE glpi_softwares.is_deleted=0`)
	_, err = glpidb.GetDB().Select(&softwares, sql)
	if err != nil {
		return softwares, err
	}
	admins, err1 := m.GetSoftwaresAdmins()
	if err != nil {
		return softwares, err1
	}
	softAdmins := []string{}
	for i, soft := range softwares {
		for _, admin := range admins {
			if soft.Groups_id_tech == admin.Id {
				softAdmins = append(softAdmins, admin.Name)
			}
		}
		if len(softAdmins) > 0 {
			softwares[i].Admins = softAdmins
			softAdmins = nil
		} else {
			softwares[i].Admins = []string{}
		}
	}
	return softwares, err
}
func (m GLPIModel) GetSoftwaresAdmins() (admins []SoftwareAdmins, err error) {
	sql := fmt.Sprintf(
		`SELECT  glpi_groups_users.groups_id AS 'id', glpi_users.name AS 'name' FROM glpi_groups_users
		 LEFT JOIN glpi_users ON glpi_users.id=glpi_groups_users.users_id
		 WHERE glpi_groups_users.groups_id IN (SELECT DISTINCT glpi_softwares.groups_id_tech FROM glpi_softwares)`)
	_, err = glpidb.GetDB().Select(&admins, sql)
	return admins, err
}

func (m GLPIModel) GetUsers() (users []GLPIUserShort, err error) {
	sql := fmt.Sprintf(
		`SELECT  u.id as id, u.name as name, IFNULL(IF(realname='' AND firstname='', u.name , CONCAT(realname, ' ',firstname)),'-') AS realname,
		authtype, LOWER(IFNULL(glpi_authldaps.name,'-')) AS ad, IFNULL(last_login, '-') as last_login,
		(SELECT COUNT(glpi_tickets_users.id) FROM glpi_tickets_users INNER JOIN glpi_tickets ON glpi_tickets.id=glpi_tickets_users.tickets_id WHERE glpi_tickets_users.users_id=u.id AND glpi_tickets_users.type=1 AND is_deleted=0) AS author
		FROM glpi_users u LEFT JOIN glpi_authldaps ON glpi_authldaps.id=u.auths_id WHERE is_deleted=false `)
	_, err = glpidb.GetDB().Select(&users, sql)
	return users, err
}

func (m GLPIModel) GetStatTickets() (tickets []TicketsStats, err error) {
	sql := fmt.Sprintf(
		`SELECT COUNT(id) AS count, TYPE AS type, YEAR (DATE) AS year, MONTH (DATE)  AS month FROM glpi_tickets
		WHERE YEAR (DATE)>2020 GROUP BY MONTH (date) , YEAR (DATE), TYPE ORDER BY DATE`)
	_, err = glpidb.GetDB().Select(&tickets, sql)
	return tickets, err
}
func (m GLPIModel) GetStatFailures() (tickets []TicketsStats, err error) {
	sql := fmt.Sprintf(
		`SELECT COUNT(id) AS count, YEAR (DATE) AS year, MONTH (DATE)  AS month  FROM (SELECT glpi_tickets.id, glpi_tickets.date from glpi_tickets
			LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			INNER JOIN glpi_plugin_fields_ticketfailures ON glpi_plugin_fields_ticketfailures.items_id=glpi_tickets.id
			INNER JOIN ( select id from glpi_plugin_fields_failcategoryfielddropdowns WHERE id>4) gpf ON gpf.id=glpi_plugin_fields_ticketfailures.plugin_fields_failcategoryfielddropdowns_id
			WHERE  glpi_tickets.is_deleted<>TRUE AND  (glpi_tickets.name not like '%%тест%%' and glpi_tickets.name not like '%%test%%')) d1
			WHERE YEAR (DATE)>2021
			GROUP BY MONTH (date) , YEAR (date)
			ORDER BY date
		 `)
	_, err = glpidb.GetDB().Select(&tickets, sql)
	return tickets, err
}

func (m GLPIModel) GetStatRegions(date string) (tickets []RegionsStats, err error) {
	sql := fmt.Sprintf(
		`SELECT count, org,ROUND(100* count/(
			SELECT count(glpi_tickets.id) FROM glpi_tickets WHERE glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.name NOT LIKE '%%ТЕСТ%%' AND glpi_tickets.name NOT LIKE '%%test%%' AND glpi_tickets.date> '%[1]s'
			)) AS proc  FROM (
			SELECT count(glpi_tickets.id) AS COUNT, 'ИРЦ Минск' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИРЦ%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.name NOT LIKE '%%ТЕСТ%%' AND glpi_tickets.name NOT LIKE '%%test%%' AND glpi_tickets.date> '%[1]s'
			UNION
			SELECT count(glpi_tickets.id) AS count, 'ИВЦ2' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИВЦ2%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.name NOT LIKE '%%ТЕСТ%%' AND glpi_tickets.name NOT LIKE '%%test%%' AND glpi_tickets.date> '%[1]s'
			UNION
			SELECT count(glpi_tickets.id) AS count, 'ИВЦ3' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИВЦ3%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.name NOT LIKE '%%ТЕСТ%%' AND glpi_tickets.name NOT LIKE '%%test%%' AND glpi_tickets.date> '%[1]s'
			UNION
			SELECT count(glpi_tickets.id) AS count, 'ИВЦ4' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИВЦ4%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.name NOT LIKE '%%ТЕСТ%%' AND glpi_tickets.name NOT LIKE '%%test%%' AND glpi_tickets.date> '%[1]s'
			UNION
			SELECT count(glpi_tickets.id) AS count, 'ИВЦ5' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИВЦ5%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.name NOT LIKE '%%ТЕСТ%%' AND glpi_tickets.name NOT LIKE '%%test%%' AND glpi_tickets.date> '%[1]s'
			UNION
			SELECT count(glpi_tickets.id) AS count, 'ИВЦ6' AS org FROM glpi_tickets LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			WHERE glpi_entities.completename LIKE '%%БЖД > ИВЦ6%%' AND glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.name NOT LIKE '%%ТЕСТ%%' AND glpi_tickets.name NOT LIKE '%%test%%' AND glpi_tickets.date> '%[1]s') a1
		 `, date)
	_, err = glpidb.GetDB().Select(&tickets, sql)
	return tickets, err
}
