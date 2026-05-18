package usecase

import (
	"encoding/json"
	"sort"
)

func (u *UseCase) FillRedisCaсheFromIUTM() error {
	categories := u.iutm.List()
	var blackList []string
	var whiteList []string
	var whiteList2 []string // если нужен третий список, уточните название

	for _, cat := range categories {
		switch cat.Name {
		case "BlackList":
			blackList = cat.Urls
		case "WhiteList":
			whiteList = cat.Urls
		case "Белый Список":
			whiteList2 = cat.Urls
			// case "ThirdList": // раскомментируйте и укажите нужное имя категории
			//     thirdList = cat.Urls
		}
	}

	// Сортировка списков
	sort.Strings(blackList)
	sort.Strings(whiteList)
	sort.Strings(whiteList2)

	// Очистка старых ключей
	u.redis.DelKeyField("utm", "blist")
	u.redis.DelKeyField("utm", "wlist")
	u.redis.DelKeyField("utm", "wlist2")
	// u.redis.DelKeyField("utm", "tlist") // для третьего списка

	// Сохранение в Redis
	if jsonList, err := json.Marshal(blackList); err == nil {
		u.redis.AddKeyFieldValue("utm", "blist", jsonList)
	}

	if jsonList, err := json.Marshal(whiteList); err == nil {
		u.redis.AddKeyFieldValue("utm", "wlist", jsonList)
	}

	if jsonList, err := json.Marshal(whiteList2); err == nil && len(whiteList2) > 0 {
		u.redis.AddKeyFieldValue("utm", "wlist2", jsonList)
	}

	return nil
}
