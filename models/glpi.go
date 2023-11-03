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
	Profiles []GLPIUserProfile `json:"profiles"`
	Groups   []IdName          `json:"groups"`
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
type IdNameType struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Type int64  `db:"type" json:"type"`
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
type Otkaz struct {
	Id            int64  `db:"id" json:"id"`
	Krit          int64  `db:"krit" json:"krit"`
	Category      int64  `db:"category" json:"category"`
	Name          string `db:"name" json:"name"`
	Status        int64  `db:"status" json:"status"`
	Impact        int64  `db:"impact" json:"impact"`
	Company       string `db:"company" json:"company"`
	Date          string `db:"date" json:"date"`
	SolvedateDate string `db:"solvedate" json:"solvedate"`
	Problems      string `db:"problems" json:"problems"`
	Content       string `db:"content" json:"content"`
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

type Ticket struct {
	Id           string `db:"id" json:"id"`
	Category     string `db:"category" json:"category"`
	Status       string `db:"status" json:"status"`
	Deleted      int64  `db:"is_deleted" json:"is_deleted"`
	Type         int64  `db:"type" json:"type"`
	RequestType  int64  `db:"requesttypes_id" json:"requesttypes_id"`
	Impact       string `db:"impact" json:"impact"`
	Date         string `db:"date" json:"date"`
	DateMod      string `db:"date_mod" json:"date_mod"`
	DateCreation string `db:"date_creation" json:"date_creation"`
	SolveDate    string `db:"solvedate" json:"solvedate"`
	Closedate    string `db:"closedate" json:"closedate,omitempty"`
	Name         string `db:"name" json:"name"`
	Content      string `db:"content" json:"content"`
	Author       string `db:"author" json:"author"`
	Orgs         string `db:"orgs" json:"orgs"`
	Company      string `db:"company" json:"company"`
	Eid          int64  `db:"eid" json:"eid"`
	Users        string `db:"users" json:"users"`
	UserGroups   string `db:"user_groups" json:"user_groups"`
	MyTicket     int64  `db:"my_ticket" json:"my_ticket"`
	Works        []Work `json:"works"`
}
type Work struct {
	Id           string `db:"id" json:"id"`
	Content      string `db:"content" json:"content"`
	DateMod      string `db:"date_mod" json:"date_mod"`
	DateCreation string `db:"date_creation" json:"date_creation"`
	Name         string `db:"name" json:"name"`
	Author       string `db:"author" json:"author"`
	IsPrivate    string `db:"is_private" json:"is_private"`
	Type         string `db:"type" json:"type"`
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
func (m GLPIModel) GetUserGroups(id int64) (groups []IdName, err error) {
	sql := fmt.Sprintf(
		`SELECT glpi_groups_users.groups_id AS id, glpi_groups.name AS name
		from glpi_groups_users INNER JOIN glpi_groups ON glpi_groups.id=glpi_groups_users.groups_id
		WHERE users_id=%d`, id)
	_, err = glpidb.GetDB().Select(&groups, sql)
	return groups, err
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

// ************************* Все отказы **********************************//
func (m GLPIModel) GetOtkazes(startdate string, enddate string) (otkazes []Otkaz, err error) {

	sql := fmt.Sprintf(
		`SELECT glpi_tickets.id as 'id',
			CASE
				WHEN STATUS<5 AND glpi_plugin_fields_failcategoryfielddropdowns.id>8
					THEN 3
				WHEN STATUS>=5 AND glpi_plugin_fields_failcategoryfielddropdowns.id>8
					THEN 2
				WHEN STATUS<5 AND (glpi_plugin_fields_failcategoryfielddropdowns.id<9 or ISNULL(glpi_plugin_fields_failcategoryfielddropdowns.id))
					THEN 1
				ELSE 0
			END AS 'krit',
			glpi_plugin_fields_failcategoryfielddropdowns.id as 'category',
			glpi_tickets.name AS 'name',
			glpi_tickets.status
				AS 'status',
			glpi_tickets.impact
				AS 'impact',
			glpi_entities.completename as 'company',

			glpi_tickets.date AS 'date',
			IFNULL(glpi_tickets.solvedate, '') AS 'solvedate',
			IFNULL((SELECT CONCAT('[',GROUP_CONCAT(CONCAT('{"id":', glpi_problems_tickets.id, ', "name":',glpi_problems_tickets.problems_id,'}') ),']') as countpr FROM glpi_problems_tickets WHERE glpi_problems_tickets.tickets_id= glpi_tickets.id GROUP BY glpi_problems_tickets.tickets_id),'[]') AS 'problems',
			glpi_tickets.content AS 'content'
			FROM glpi_tickets
			LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
			LEFT JOIN glpi_plugin_fields_ticketfailures ON glpi_plugin_fields_ticketfailures.items_id=glpi_tickets.id
			LEFT JOIN glpi_plugin_fields_failcategoryfielddropdowns ON glpi_plugin_fields_failcategoryfielddropdowns.id=glpi_plugin_fields_ticketfailures.plugin_fields_failcategoryfielddropdowns_id
			WHERE glpi_tickets.is_deleted<>TRUE  AND glpi_plugin_fields_failcategoryfielddropdowns.id>4
			AND glpi_tickets.name not like '%%тест%%' and glpi_tickets.name not like '%%test%%' AND
			((date>='%[1]s' AND date <='%[2]s') OR (solvedate>='%[1]s' AND solvedate <='%[2]s') OR (date<'%[1]s' AND solvedate >'%[2]s') OR (date<'%[1]s' AND solvedate is null))
			ORDER BY date desc
		`, startdate, enddate)
	_, err = glpidb.GetDB().Select(&otkazes, sql)
	return otkazes, err
}

// ************************* Всего отказов **********************************//
func (m GLPIModel) GetStatOtkazSum() (sum int64, err error) {
	sum = 0
	sql := fmt.Sprintf(
		`SELECT count(glpi_tickets.id) AS "sum" FROM glpi_tickets
		LEFT JOIN glpi_entities ON glpi_tickets.entities_id = glpi_entities.id
		INNER JOIN glpi_plugin_fields_ticketfailures ON glpi_plugin_fields_ticketfailures.items_id=glpi_tickets.id
		INNER JOIN ( select id from glpi_plugin_fields_failcategoryfielddropdowns WHERE id>4) gpf ON gpf.id=glpi_plugin_fields_ticketfailures.plugin_fields_failcategoryfielddropdowns_id
		WHERE  glpi_tickets.is_deleted<>TRUE  AND glpi_tickets.name not like "%%тест%%" and glpi_tickets.name not like "%%test%%"`)
	err = glpidb.GetDB().QueryRow(sql).Scan(&sum)
	return sum, err
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

// ************************* Незакрытые заявки **********************************//

func (m GLPIModel) GetTicketsNonClosed() (tickets []Ticket, err error) {

	sql := fmt.Sprintf(
		`
		SELECT IFNULL(JSON_EXTRACT(ancestors_cache, '$.*'),JSON_ARRAY(0)) AS 'orgs',IFNULL(fkat,0) as category,
		gt.id , gt.content, gt.status, gt.name, gt.impact, glpi_entities.completename as company, glpi_entities.id as eid, IFNULL(gt.date,'') as date, gt.date_mod, gt.date_creation, IFNULL(gt.solvedate,'') as solvedate,
		CONCAT(ifnull(NULLIF(glpi_users.realname, ''), 'не опреденен'),' ', ifnull(NULLIF(glpi_users.firstname, ''),'')) AS author,
		IFNULL((SELECT CONCAT('[',GROUP_CONCAT(CONCAT('{"id":', glpi_tickets_users.users_id, ', "name":"',glpi_users.name,'"', ', "type":',glpi_tickets_users.type,'}') ),']')  FROM glpi_tickets_users INNER JOIN glpi_users ON glpi_users.id=glpi_tickets_users.users_id  WHERE glpi_tickets_users.tickets_id = gt.id  ),'[]') AS users,
		IFNULL((SELECT CONCAT('[',GROUP_CONCAT(CONCAT('{"id":', glpi_groups_tickets.groups_id, ', "name":"',glpi_groups.name,'"', ', "type":',glpi_groups_tickets.type,'}') ),']')  FROM glpi_groups_tickets INNER JOIN glpi_groups ON glpi_groups.id=glpi_groups_tickets.groups_id  WHERE glpi_groups_tickets.tickets_id = gt.id  ),'[]') AS user_groups
		FROM (SELECT * FROM glpi_tickets WHERE glpi_tickets.is_deleted=0  AND STATUS <5 AND LOWER(glpi_tickets.name) not like '%%т1е1с1т%%' AND LOWER(glpi_tickets.name) not like '%%test%%') gt
		INNER JOIN glpi_entities ON gt.entities_id=glpi_entities.id
		LEFT JOIN glpi_users ON gt.users_id_recipient=glpi_users.id
		LEFT JOIN  (SELECT items_id,plugin_fields_failcategoryfielddropdowns_id AS fkat  from glpi_plugin_fields_ticketfailures WHERE plugin_fields_failcategoryfielddropdowns_id>4) fc ON fc.items_id=gt.id
		`)
	_, err = glpidb.GetDB().Select(&tickets, sql)
	return tickets, err
}

func (m GLPIModel) GetTicket(id string) (ticket Ticket, err error) {

	sql := fmt.Sprintf(
		`
		SELECT IFNULL(JSON_EXTRACT(ancestors_cache, '$.*'),JSON_ARRAY(0)) AS 'orgs',IFNULL(fkat,0) as category,
		gt.id , gt.content, gt.status, gt.name, gt.impact, glpi_entities.completename as company, glpi_entities.id as eid,
		IFNULL(gt.date,'') as date, gt.date_mod, IFNULL(gt.closedate,'') as closedate, gt.date_creation, IFNULL(gt.solvedate,'') as solvedate,
		gt.is_deleted, gt.type, gt.requesttypes_id,
		CONCAT(ifnull(NULLIF(glpi_users.realname, ''), 'не опреденен'),' ', ifnull(NULLIF(glpi_users.firstname, ''),'')) AS author,
		IFNULL((SELECT CONCAT('[',GROUP_CONCAT(CONCAT('{"id":', glpi_tickets_users.users_id, ', "name":"',glpi_users.name,'"', ', "type":',glpi_tickets_users.type,'}') ),']')  FROM glpi_tickets_users INNER JOIN glpi_users ON glpi_users.id=glpi_tickets_users.users_id  WHERE glpi_tickets_users.tickets_id = gt.id  ),'[]') AS users,
		IFNULL((SELECT CONCAT('[',GROUP_CONCAT(CONCAT('{"id":', glpi_groups_tickets.groups_id, ', "name":"',glpi_groups.name,'"', ', "type":',glpi_groups_tickets.type,'}') ),']')  FROM glpi_groups_tickets INNER JOIN glpi_groups ON glpi_groups.id=glpi_groups_tickets.groups_id  WHERE glpi_groups_tickets.tickets_id = gt.id  ),'[]') AS user_groups
		FROM (SELECT * FROM glpi_tickets WHERE glpi_tickets.id=%s) gt
		INNER JOIN glpi_entities ON gt.entities_id=glpi_entities.id
		LEFT JOIN glpi_users ON gt.users_id_recipient=glpi_users.id
		LEFT JOIN  (SELECT items_id,plugin_fields_failcategoryfielddropdowns_id AS fkat  from glpi_plugin_fields_ticketfailures WHERE plugin_fields_failcategoryfielddropdowns_id>4) fc ON fc.items_id=gt.id
		`, id)
	err = glpidb.GetDB().SelectOne(&ticket, sql)

	return ticket, err

}

func (m GLPIModel) TicketWorks(ticketID string) (work []Work, err error) {
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
	  from glpi_tickets WHERE id=%[1]s
							`, ticketID)
	_, err = glpidb.GetDB().Select(&work, proc)
	if err != nil {
		return nil, err
	}
	return work, nil
}
