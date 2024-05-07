package entity

type GlpiApiToken struct {
	Token string `json:"session_token"`
}

type GLPIApiResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type NewTicketInputForm struct {
	Input NewTicketForm `form:"input" json:"input"`
}

type UsersIdRequesterNotifStruct struct {
	UseNotification  string   `form:"use_notification" json:"use_notification" `
	AlternativeEmail []string `form:"alternative_email" json:"alternative_email" `
}

type NewTicketForm struct {
	Name                  string                      `form:"name" json:"name" binding:"required"`
	Content               string                      `form:"content" json:"content" binding:"required"`
	RequesttypesId        int                         `form:"requesttypes_id" json:"requesttypes_id" binding:"required"`
	UsersIdRequester      int                         `form:"_users_id_requester" json:"_users_id_requester" binding:"required"`
	UsersIdRequesterNotif UsersIdRequesterNotifStruct `form:"_users_id_requester_notif" json:"_users_id_requester_notif" `
	EntitiesId            int                         `form:"entities_id" json:"entities_id" `
	Type                  int                         `form:"type" json:"type" binding:"required"`
	Urgency               int                         `form:"urgency" json:"urgency" binding:"required"`
	UsersIdAssign         int                         `form:"_users_id_assign" json:"_users_id_assign"`
	GroupsIdAssign        int                         `form:"_groups_id_assign" json:"_groups_id_assign"`
	Status                int                         `form:"status" json:"status"`
	TicketId              int                         `form:"tickets_id" json:"tickets_id"`
	User                  string                      `form:"user" json:"user"`
	Token                 string                      `form:"token" json:"token"`
}

type UpdateTicketForm struct {
	TicketId       string `form:"tickets_id" json:"tickets_id" binding:"required"`
	Name           string `form:"name" json:"name,omitempty" binding:"max=100"`
	Content        string `form:"content" json:"content,omitempty" binding:"max=400"`
	Urgency        string `form:"urgency" json:"urgency,omitempty" `
	UsersIdAssign  string `form:"_users_id_assign" json:"_users_id_assign,omitempty"`
	GroupsIdAssign string `form:"_groups_id_assign" json:"_groups_id_assign,omitempty"`
	Status         string `form:"status" json:"status,omitempty"`
}

// GLPI API структуры комментария заявки

type NewCommentInputForm struct {
	Input NewCommentForm `form:"input" json:"input"`
}
type NewCommentForm struct {
	ItemType        string `form:"itemtype" json:"itemtype"`
	ItemId          int    `form:"items_id" json:"items_id" binding:"required"`
	Content         string `form:"content" json:"content" binding:"required,max=1000"`
	RequestTypesId  int    `form:"requesttypes_id" json:"requesttypes_id"`
	IsPrivate       bool   `form:"is_private" json:"is_private"`
	Status          int    `form:"status" json:"status"`
	SolutiontypesId int    `form:"solutiontypes_id" json:"solutiontypes_id"`
	User            string `form:"user" json:"user"`
	Token           string `form:"token" json:"token"`
}

type GLPITicketUserInputForm struct {
	Input GLPITicketUserForm `form:"input" json:"input"`
}
type GLPITicketUserForm struct {
	TicketId int    `form:"tickets_id" json:"tickets_id" binding:"required"`
	UsersId  int    `form:"users_id" json:"users_id" binding:"required"`
	Type     int    `form:"type" json:"type" binding:"required"`
	User     string `form:"user" json:"user"`
	Token    string `form:"token" json:"token"`
}
