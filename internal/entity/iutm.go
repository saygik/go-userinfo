package entity

type IutmCategoryList struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Comment string   `json:"comment"`
	Urls    []string `json:"urls"`
}
