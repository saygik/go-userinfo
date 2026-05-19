package usecase

import (
	"encoding/json"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) IUTMGetList(list string) ([]string, error) {
	var r []string
	redisLists, err := u.redis.GetKeyFieldValue("utm", list) //"wlist"
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(redisLists), &r)
	return r, nil
}

func (u *UseCase) IUTMGetWlist2() ([]string, error) {
	var r []string

	redisLists, err := u.redis.GetKeyFieldValue("utm", "wlist2")
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(redisLists), &r)

	return r, nil
}

func (u *UseCase) IUTMGetBlist() ([]string, error) {
	var r []string

	redisLists, err := u.redis.GetKeyFieldValue("utm", "blist")
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(redisLists), &r)

	return r, nil
}

func (u *UseCase) IUTMGetAllLists() ([]entity.UrlInfo, error) {
	var result []entity.UrlInfo

	lists := []struct {
		Key  string
		Name string
	}{
		{Key: "wlist", Name: "Белый список ЦЗИ"},
		{Key: "blist", Name: "Черный список ЦЗИ"},
		{Key: "wlist2", Name: "Белый список НОД2"},
	}

	for _, list := range lists {
		redisValue, err := u.redis.GetKeyFieldValue("utm", list.Key)
		if err != nil {
			continue
		}

		var urls []string
		if err := json.Unmarshal([]byte(redisValue), &urls); err != nil {
			continue
		}

		// Каждый URL отдельным элементом
		for _, url := range urls {
			result = append(result, entity.UrlInfo{
				Name: list.Name,
				Url:  url,
			})
		}
	}

	return result, nil
}
