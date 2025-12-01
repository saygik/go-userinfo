package usecase

import (
	"encoding/json"
	"sort"
)

func (u *UseCase) FillRedisCaсheFromIUTM() error {
	categories := u.iutm.List()
	var urls []string
	for _, cat := range categories {
		if cat.Name == "WhiteList" {
			urls = cat.Urls
			break
		}
	}
	sort.Strings(urls)
	u.redis.DelKeyField("utm", "wlist")
	jsonList, _ := json.Marshal(urls)
	u.redis.AddKeyFieldValue("utm", "wlist", jsonList)
	_ = urls
	return nil
}
