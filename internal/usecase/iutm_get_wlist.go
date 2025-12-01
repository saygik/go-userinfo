package usecase

import "encoding/json"

func (u *UseCase) IUTMGetWlist() ([]string, error) {
	var r []string

	redisLists, err := u.redis.GetKeyFieldValue("utm", "wlist")
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(redisLists), &r)

	return r, nil
}
