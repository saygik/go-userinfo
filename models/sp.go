package models

import (
	"encoding/json"
	"github.com/saygik/go-userinfo/sp"
)

//Sharepoint Model ...
type SPModel struct{}

type urlStruct struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}
type Zal struct {
	Id     int       `json:"id"`
	Order  int       `json:"order0"`
	ConfId int       `json:"confId"`
	Zal    urlStruct `json:"zal"`
}

//GLPI User find by Mail ...
func (m SPModel) All() (zals []Zal, err error) {
	list := sp.GetSP().Web().GetList("Lists/zals")
	data, err := list.Items().Select("Id,order0,ConfId,zal").Get()
	if err != nil {
		return zals, err
	}
	if err := json.Unmarshal(data.Normalized(), &zals); err != nil {
		return zals, err
	}
	return zals, err
}
