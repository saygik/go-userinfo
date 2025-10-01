package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/saygik/go-userinfo/internal/state"
)

func (u *UseCase) FillRedisCaсheFromAD() error {
	adl := u.ad.DomainList()
	//	u.redis.ClearAllDomainsUsers()
	u.redis.Delete("allusers")
	for _, one := range adl {
		u.log.Info("Получение данных домена " + one.Name + "...")
		users, err := u.ad.GetDomainUsers(one.Name)
		comps, _ := u.ad.GetDomainComputers(one.Name)
		rmsPort := u.ad.GetDomainRMSPort(one.Name)
		u.log.Info("Получение данных домена " + one.Name + " завершено. Обработка данных...")
		if err == nil || len(users) > 0 {
			println("Get from ad to redis from " + one.Name)
			ips, _ := u.repo.GetDomainUsersIP(one.Name)
			avatars, _ := u.repo.GetDomainUsersAvatars(one.Name)
			for _, user := range users {

				user["domain"] = one
				user["rms_port"] = rmsPort
				if val, ok := user["userAccountControl"]; ok {
					uac := val.(string)
					if uac == "514" || uac == "66050" {
						user["disabled"] = true
					}
				}

				if val, ok := user["lockoutTime"]; ok {
					if val.(string) != "0" {
						lockoutTime, err := ADFiletimeToGoTime(val.(string))
						if err == nil {
							user["locked"] = true
							user["lockoutTime"] = lockoutTime

						}
					}
				}

				arch := user["distinguishedName"].(string)
				arch = strings.ToUpper(arch)
				if strings.Contains(arch, "=АРХИВ") {
					user["archive"] = true
				}
				if IsStringInArray("Отключенные Кадровичком", user["memberOf"]) {
					user["sap_disabled"] = true
				}
				if IsStringInArray("Пользователи интернета", user["memberOf"]) {
					user["internet"] = true
				}
				if IsStringInArray("Пользователи интернета Белый список", user["memberOf"]) {
					user["internetwl"] = true
				}
				if len(ips) > 0 {
					for _, ip := range ips {
						if isStringObjsEqual(user["userPrincipalName"], ip.Login) {
							user["ip"] = ip.Ip
							user["computer"] = ip.Computer
							user["ip_date"] = ip.IpDate
							user["rms_installed"] = ip.Rms
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
				upn, ok := user["userPrincipalName"].(string)
				if ok {
					u.redis.AddKeyFieldValue("allusers", upn, jsonUser)
				}

				//				redisClient.HSet(ctx, "allusers", user["userPrincipalName"], jsonUser).Err()
			}

		}
		sort.Slice(users, func(i, j int) bool {
			return fmt.Sprintf("%v", users[i]["cn"]) < fmt.Sprintf("%v", users[j]["cn"])
		})
		jsonUsers, _ := json.Marshal(users)
		jsonComps, _ := json.Marshal(comps)
		u.redis.DelKeyField("ad", one.Name)
		err1 := u.redis.AddKeyFieldValue("ad", one.Name, jsonUsers)
		if err1 != nil {
			return errors.New("key does not exists")
		}
		u.redis.DelKeyField("adc", one.Name)
		u.redis.AddKeyFieldValue("adc", one.Name, jsonComps)
		u.log.Info("Получение данных домена " + one.Name + " завершено. Обработка данных завершена.")
	}
	u.log.Info("Получение данных всех доменов завершено.")
	if !state.IsInitialized() {
		state.SetInitialized()
	}
	return nil
}
