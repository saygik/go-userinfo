package entity

type GlpiApiToken struct {
	Token string `json:"session_token"`
}

type GLPIApiResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
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
