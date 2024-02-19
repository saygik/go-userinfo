package usecase

func (u *UseCase) GetADComputers(user string) ([]map[string]interface{}, error) {

	var res []map[string]interface{}
	redisADComputers, err := u.redis.GetKeyFieldAll("adc")
	if err != nil {
		return nil, err
	}
	for domainName, oneDomain := range redisADComputers {
		access := u.GetAccessToResource(domainName, user)
		if access == -1 {
			continue
		}
		var r []map[string]interface{}
		unmarshalString(oneDomain, &r)

		res = append(res, r...)
	}
	return res, nil
}
