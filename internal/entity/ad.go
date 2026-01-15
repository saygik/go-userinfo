package entity

type DomainList struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type AvatarForm struct {
	Avatar string `form:"avatar" json:"avatar"`
}

type ADInternetGroups struct {
	WhiteList []string `json:"whitelist"`
	Full      []string `json:"full"`
	Tech      []string `json:"tech"`
}

type ADInternetGroupsDN struct {
	WhiteList string `json:"whitelist"`
	Full      string `json:"full"`
	Tech      string `json:"tech"`
}
