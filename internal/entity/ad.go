package entity

type DomainList struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type AvatarForm struct {
	Avatar string `form:"avatar" json:"avatar"`
}
