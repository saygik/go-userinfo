package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

func (u *UseCase) ClearRedisCaсhe() {
	//u.r.clearAllDomainsUsers()
}

func (u *UseCase) FillRedisCaсheFromAD() error {
	adl := u.ad.DomainList()
	for _, one := range adl {
		users, err := u.ad.GetDomainUsers("brnv.rw")
		comps, _ := u.ad.GetDomainComputers("brnv.rw")
		if err == nil || len(users) > 0 {
			println("Get from ad to redis from " + one)
			ips, _ := u.repo.GetDomainUsersIP(one)
			avatars, _ := u.repo.GetDomainUsersAvatars(one)
			for _, user := range users {
				user["domain"] = one
				if IsStringInArray("Пользователи интернета", user["memberOf"]) {
					user["internet"] = true
				}
				if IsStringInArray("Пользователи интернета Белый список", user["memberOf"]) {
					user["internetwl"] = true
				}
				if len(ips) > 0 {
					for _, ip := range ips {
						if user["userPrincipalName"] == ip.Login {
							user["ip"] = ip.Ip
							user["computer"] = ip.Computer
						}
					}
				}
				if len(avatars) > 0 {
					for _, avatar := range avatars {
						if user["userPrincipalName"] == avatar.Name {
							user["avatar"] = avatar.Avatar
						}
					}
				}

				jsonUser, _ := json.Marshal(user)
				u.redis.AddKeyFieldValue("allusers", user["userPrincipalName"].(string), jsonUser)
				//				redisClient.HSet(ctx, "allusers", user["userPrincipalName"], jsonUser).Err()
			}

		}
		sort.Slice(users, func(i, j int) bool {
			return fmt.Sprintf("%v", users[i]["cn"]) < fmt.Sprintf("%v", users[j]["cn"])
		})
		jsonUsers, _ := json.Marshal(users)
		jsonComps, _ := json.Marshal(comps)
		err1 := u.redis.AddKeyFieldValue("ad", one, jsonUsers)
		if err1 != nil {
			return errors.New("key does not exists")
		}
		u.redis.AddKeyFieldValue("adc", one, jsonComps)
	}
	return nil
}
