package entity

type MattermostUser struct {
	Id          string `json:"id"`
	Name        string `json:"username"`
	AuthService string `json:"auth_service,omitempty"`
	AD          string `json:"ad,omitempty"`
	Email       string `json:"email"`
	Nickname    string `json:"nickname,omitempty"`
	Roles       string `json:"roles,omitempty"`
	IsBot       bool   `json:"is_bot,omitempty"`
	Registred   bool   `json:"registred"`
}
