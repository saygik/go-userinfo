package usecase

import (
	"encoding/json"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) GetADUsersPublicInfo(perm entity.Permissions) ([]map[string]interface{}, error) {

	if !perm.IsTech {
		return nil, u.Error("у вас нет прав на просмотр списка пользователей всех доменов")
	}
	redisADUsers, err := u.redis.GetKeyFieldAll("allusers")
	if err != nil {
		return nil, err
	}
	var users []map[string]interface{}
	for _, value := range redisADUsers {
		var user map[string]interface{}
		json.Unmarshal([]byte(value), &user)
		for _, field := range sensitiveLongFields {
			delete(user, field)
		}
		user["findedInAD"] = true
		user["name"] = user["userPrincipalName"]
		users = append(users, user)
	}
	//	json.Unmarshal([]byte(redisADUsers), &users)
	return users, nil
}
