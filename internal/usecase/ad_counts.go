package usecase

import "encoding/json"

func (u *UseCase) GetAdCounts() (int, int, error) {
	var r []map[string]interface{}

	redisADUsers, err := u.redis.GetKeyFieldAll("ad")
	if err != nil {
		return 0, 0, err
	}
	redisADComputers, err := u.redis.GetKeyFieldAll("adc")
	if err != nil {
		return 0, 0, err
	}
	users := 0
	computers := 0
	for _, oneDomain := range redisADUsers {
		json.Unmarshal([]byte(oneDomain), &r)
		users = users + len(r)
	}
	for _, oneDomain := range redisADComputers {
		json.Unmarshal([]byte(oneDomain), &r)
		computers = computers + len(r)
	}

	return users, computers, nil
}
