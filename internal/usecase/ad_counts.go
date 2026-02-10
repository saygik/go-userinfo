package usecase

import "encoding/json"

// GetAdCounts возвращает общее количество пользователей и компьютеров по всем доменам.
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
		_ = json.Unmarshal([]byte(oneDomain), &r)
		users += len(r)
	}
	for _, oneDomain := range redisADComputers {
		_ = json.Unmarshal([]byte(oneDomain), &r)
		computers += len(r)
	}

	return users, computers, nil
}

// GetAdCountsDomain возвращает количество пользователей и компьютеров только для одного домена.
func (u *UseCase) GetAdCountsDomain(domain string) (int, int, error) {
	var usersArr []map[string]interface{}
	var compsArr []map[string]interface{}

	oneDomainUsers, err := u.redis.GetKeyFieldValue("ad", domain)
	if err != nil {
		return 0, 0, err
	}
	if err := json.Unmarshal([]byte(oneDomainUsers), &usersArr); err != nil {
		return 0, 0, err
	}

	oneDomainComps, err := u.redis.GetKeyFieldValue("adc", domain)
	if err != nil {
		return 0, 0, err
	}
	if err := json.Unmarshal([]byte(oneDomainComps), &compsArr); err != nil {
		return 0, 0, err
	}

	return len(usersArr), len(compsArr), nil
}
