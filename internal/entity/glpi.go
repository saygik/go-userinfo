package entity

type GLPIUserProfile struct {
	Id        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	EName     string `db:"ename" json:"ename"`
	Eid       int    `db:"eid" json:"eid"`
	Orgs      string `db:"orgs" json:"orgs"`
	Recursive bool   `db:"recursive" json:"recursive"`
}

type GLPIGroup struct {
	Id       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Type     int    `db:"type" json:"type"`
	Presence bool   `db:"presence" json:"presence"`
}

type GLPIUser struct {
	Id             int               `db:"id" json:"id"`
	Name           string            `db:"name" json:"name"`
	Fio            string            `db:"fio" json:"fio"`
	Email          string            `db:"email" json:"email"`
	Self           string            `db:"self" json:"self"`
	SelfId         int               `db:"selfid" json:"selfid"`
	Date           string            `db:"date" json:"date"`
	Org            string            `db:"org" json:"org"`
	Api            bool              `db:"api" json:"api"`
	Type           int               `db:"type" json:"type"`
	Author         int               `db:"author" json:"author"`
	Executor       int               `db:"executor" json:"executor"`
	Profiles       []GLPIUserProfile `json:"profiles"`
	Groups         []IdName          `json:"groups"`
	TrackingGroups []IdName          `json:"tracking_groups"`
	Tickets        []GLPI_Ticket     `json:"tickets"`
}

type GLPIUserShort struct {
	Id        int                    `db:"id" json:"id"`
	Name      string                 `db:"name" json:"name"`
	Realname  string                 `db:"realname" json:"realname"`
	Authtype  int                    `db:"authtype" json:"authtype"`
	Deleted   bool                   `db:"is_deleted" json:"is_deleted"`
	AD        string                 `db:"ad" json:"ad"`
	LastLogin string                 `db:"last_login" json:"last_login"`
	Author    int                    `db:"author" json:"author"`
	ADProfile map[string]interface{} `json:"ad_profile"`
}

type GLPI_Work struct {
	Id           string     `db:"id" json:"id"`
	Content      string     `db:"content" json:"content"`
	DateMod      string     `db:"date_mod" json:"date_mod"`
	DateCreation string     `db:"date_creation" json:"date_creation"`
	Name         string     `db:"name" json:"name"`
	Author       string     `db:"author" json:"author"`
	AuthorProps  SimpleUser `db:"author_props" json:"author_props"`
	IsPrivate    bool       `db:"is_private" json:"is_private"`
	Type         string     `db:"type" json:"type"`
	Status       int        `db:"status" json:"status"`
}

type GLPI_Ticket_Profile struct {
	Id      string `db:"id" json:"id"`
	Content string `db:"content" json:"content"`
	Name    string `db:"name" json:"name"`
	Company string `db:"company" json:"company"`
}

// GLPI_Ticket присутствие пользователя в заявке
type GLPI_Ticket_User_Presence struct {
	Initiator bool ` json:"initiator"`
	Executor  bool ` json:"executor"`
	Observer  bool ` json:"observer"`
}

// GLPI_Ticket присутствие пользователя в заявке
type GLPI_Ticket_Users struct {
	Initiators []map[string]interface{} ` json:"initiators"`
	Executors  []map[string]interface{} ` json:"executors"`
	Observers  []map[string]interface{} ` json:"observers"`
}

// GLPI_Ticket присутствие пользователя в заявке
type GLPI_Ticket_Groups struct {
	Initiators []GLPIGroup ` json:"initiators"`
	Executors  []GLPIGroup ` json:"executors"`
	Observers  []GLPIGroup ` json:"observers"`
}

// GLPI_Ticket заявка GLPI
type GLPI_Ticket struct {
	Id            int                       `db:"id" json:"id"`
	Krit          int                       `db:"krit" json:"krit"`
	Category      int                       `db:"category" json:"category"`
	Status        int                       `db:"status" json:"status"`
	Deleted       bool                      `db:"is_deleted" json:"is_deleted"`
	Type          int                       `db:"type" json:"type"`
	RequestType   int                       `db:"requesttypes_id" json:"requesttypes_id"`
	Impact        int                       `db:"impact" json:"impact"`
	Date          string                    `db:"date" json:"date"`
	DateMod       string                    `db:"date_mod" json:"date_mod"`
	DateCreation  string                    `db:"date_creation" json:"date_creation"`
	SolveDate     string                    `db:"solvedate" json:"solvedate"`
	Closedate     string                    `db:"closedate" json:"closedate,omitempty"`
	Name          string                    `db:"name" json:"name"`
	Content       string                    `db:"content" json:"content"`
	Author        string                    `db:"author" json:"author"`
	Orgs          string                    `db:"orgs" json:"orgs"`
	Company       string                    `db:"company" json:"company"`
	Eid           int                       `db:"eid" json:"eid"`
	UsersS        string                    `db:"users_s" json:"users_s,omitempty"`
	Users         GLPI_Ticket_Users         `json:"users"`
	UserGroupsS   string                    `db:"user_groups_s" json:"user_groups_s,omitempty"`
	Groups        GLPI_Ticket_Groups        `json:"groups"`
	MyTicket      int                       `db:"my_ticket" json:"my_ticket"`
	ProblemsCount int                       `db:"problemscount" json:"problemscount"`
	RecipientId   int                       `db:"recipient_id" json:"recipient_id"`
	Group         IdName                    `db:"group" json:"group"`
	UserPresence  GLPI_Ticket_User_Presence `json:"user_presence"`
	GroupPresence GLPI_Ticket_User_Presence `json:"group_presence"`
	Recipient     map[string]interface{}    `json:"recipient"`
	Works         []GLPI_Work               `json:"works"`
	Problems      []GLPI_Problem            `json:"problems"`
	GroupId       int                       `db:"group_id" json:"group_id"`
	GroupName     string                    `db:"group_name" json:"group_name"`
}

type GLPI_Otkaz struct {
	Id            int    `db:"id" json:"id"`
	Krit          int    `db:"krit" json:"krit"`
	Category      int    `db:"category" json:"category"`
	Name          string `db:"name" json:"name"`
	Deleted       bool   `db:"deleted" json:"deleted"`
	Status        int    `db:"status" json:"status"`
	Impact        int    `db:"impact" json:"impact"`
	Company       string `db:"company" json:"company"`
	Date          string `db:"date" json:"date"`
	SolvedateDate string `db:"solvedate" json:"solvedate"`
	Problems      string `db:"problems" json:"problems"`
	Content       string `db:"content" json:"content"`
}

// GLPI_Problem заявка GLPI
type GLPI_Problem struct {
	Id             int                       `db:"id" json:"id"`
	Krit           int                       `db:"krit" json:"krit"`
	Category       int                       `db:"category" json:"category"`
	Status         int                       `db:"status" json:"status"`
	Type           int                       `db:"type" json:"type"`
	Recursive      bool                      `db:"recursive" json:"recursive"`
	Name           string                    `db:"name" json:"name"`
	Content        string                    `db:"content" json:"content"`
	ImpactContent  string                    `db:"impactcontent" json:"impactcontent"`
	CauseContent   string                    `db:"causecontent" json:"causecontent"`
	SymptomContent string                    `db:"symptomcontent" json:"symptomcontent"`
	Company        string                    `db:"company" json:"company"`
	Deleted        bool                      `db:"is_deleted" json:"is_deleted"`
	Impact         int                       `db:"impact" json:"impact"`
	Date           string                    `db:"date" json:"date"`
	DateMod        string                    `db:"datemod" json:"datemod"`
	DateCreation   string                    `db:"date_creation" json:"date_creation"`
	Solvedate      string                    `db:"solvedate" json:"solvedate"`
	Solvetime      string                    `db:"solvetime" json:"solvetime"`
	Closedate      string                    `db:"closedate" json:"closedate,omitempty"`
	TicketsId      string                    `db:"ticketsid" json:"ticketsid"`
	TicketsCount   int                       `db:"ticketscount" json:"ticketscount"`
	Author         string                    `db:"author" json:"author"`
	Orgs           string                    `db:"orgs" json:"orgs"`
	Eid            int                       `db:"eid" json:"eid"`
	UserGroups     string                    `db:"user_groups" json:"user_groups"`
	Users          GLPI_Ticket_Users         `json:"users"`
	Groups         GLPI_Ticket_Groups        `json:"groups"`
	MyTicket       int                       `db:"my_ticket" json:"my_ticket"`
	Works          []GLPI_Work               `json:"works"`
	RecipientId    int                       `db:"recipient_id" json:"recipient_id"`
	Recipient      map[string]interface{}    `json:"recipient"`
	Tickets        []GLPI_Otkaz              `json:"tickets"`
	UserPresence   GLPI_Ticket_User_Presence `json:"user_presence"`
	GroupPresence  GLPI_Ticket_User_Presence `json:"group_presence"`
}

type GLPITicketsStats struct {
	Сount int `db:"count" json:"count"`
	Type  int `db:"type" json:"type"`
	Year  int `db:"year" json:"year"`
	Month int `db:"month" json:"month"`
	Day   int `db:"day" json:"day"`
}

type RegionsDayStats struct {
	Сount int    `db:"count" json:"count"`
	Org   string `db:"org" json:"org"`
	Day   int    `db:"day" json:"day"`
}

type GLPIStatsTop10 struct {
	Name         string `db:"name" json:"name"`
	Completename string `db:"completename" json:"completename"`
	Company      string `db:"company" json:"company"`
	Count        int    `db:"count" json:"count"`
}

type GLPIStatsCounts struct {
	T1  string `db:"t1" json:"t1"`
	T11 string `db:"t1_1" json:"t1_1"`
	T12 string `db:"t1_2" json:"t1_2"`
	T2  string `db:"t2" json:"t2"`
	T21 string `db:"t2_1" json:"t2_1"`
	T22 string `db:"t2_2" json:"t2_2"`
	T3  string `db:"t3" json:"t3"`
	T31 string `db:"t3_1" json:"t3_1"`
	T32 string `db:"t3_2" json:"t3_2"`
	T33 string `db:"t3_3" json:"t3_3"`
}

type GLPIRegionsStats struct {
	Сount int    `db:"count" json:"count"`
	Org   string `db:"org" json:"org"`
	Proc  int    `db:"proc" json:"proc"`
}

type TreemapData struct {
	X string `db:"x" json:"x"`
	Y int    `db:"y" json:"y"`
}

type GLPI_Ticket_Users_Simple struct {
	Initiators []SimpleUser ` json:"initiators"`
	Executors  []SimpleUser ` json:"executors"`
	Observers  []SimpleUser ` json:"observers"`
}

type GLPI_Ticket_Report struct {
	Id                  int                      `db:"id" json:"id"`
	Name                string                   `db:"name" json:"name"`
	Content             string                   `db:"content" json:"content"`
	Type                string                   `db:"type" json:"type"`
	Status              int                      `db:"status" json:"status"`
	Statustext          string                   `db:"statustext" json:"statustext"`
	Impact              int                      `db:"impact" json:"impact"`
	Impacttext          string                   `db:"impacttext" json:"impacttext"`
	Org                 string                   `db:"org" json:"org"`
	Date                string                   `db:"date_vos" json:"date_vos"`
	DateMod             string                   `db:"date_mod" json:"date_mod"`
	DateCreation        string                   `db:"date_creation" json:"date_creation"`
	SolveDate           string                   `db:"solvedate" json:"solvedate"`
	Closedate           string                   `db:"closedate" json:"closedate"`
	FailCategory        string                   `db:"fail_category" json:"fail_category"`
	FailCategoryComment string                   `db:"fail_category_comment" json:"fail_category_comment"`
	RequestType         string                   `db:"requesttype" json:"requesttype"`
	Du                  bool                     `db:"du" json:"du"`
	Mds                 bool                     `db:"mds" json:"mds"`
	Users               GLPI_Ticket_Users_Simple `json:"users"`
	Comments            []GLPI_Work              `json:"comments"`
	Solutions           []GLPI_Work              `json:"solutions"`
	ExecutorsGroup      string                   `json:"executors_group"`
	ExecutorsGroupUsers []SimpleUser             `json:"executors_group_users"`
	Objects             GLPI_Objects             `json:"objects"`
	MyTicket            int                      `json:"my_ticket"`
	Orgs                string                   `db:"orgs" json:"orgs"`
	Eid                 int                      `db:"eid" json:"eid"`
}
type GLPI_Objects struct {
	Softwares        []GLPI_Object `json:"softwares"`
	Servers          []GLPI_Object `json:"servers"`
	NetworkEquipment []GLPI_Object `json:"network_equipment"`
}
type GLPI_Object struct {
	Name     string `db:"name" json:"name"`
	Fullname string `db:"fullname" json:"fullname"`
	Group    string `db:"group" json:"group"`
	Place    string `db:"place" json:"place"`
}
