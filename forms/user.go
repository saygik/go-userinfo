package forms

type UserActivityForm struct {
	User    string `form:"user" json:"user" binding:"required,email"`
	Ip      string `form:"ip" json:"ip"`
	Activiy string `form:"activity" json:"activity"`
}
