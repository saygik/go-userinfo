package usecase

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/saygik/go-userinfo/internal/entity"
	"github.com/saygik/go-userinfo/internal/state"
)

func (u *UseCase) FillRedisCaсheFromAD() error {
	if state.IsFillingRedis() {
		u.log.Error("кеш не обновлён, попытка заполнения кеша пользователей при невыполненной предыдущей задаче...")
		return u.Error("кеш не обновлён")
	}
	defer func() {
		if !state.IsInitialized() {
			state.SetInitialized()
		}
		state.SetFillingRedis(false)
	}()

	state.SetFillingRedis(true)
	adl := u.ad.DomainList()
	//	u.redis.ClearAllDomainsUsers()
	u.redis.Delete("allusers:stagging")
	var wg sync.WaitGroup
	for _, one := range adl {
		wg.Add(1)
		go func(one entity.DomainList) {
			enabledUsersCount, disabledUsersCount, lockedUsersCount := 0, 0, 0
			fullInternetUsersCount, whiteListInternetUsersCount, techInternetUsersCount := 0, 0, 0
			start := time.Now()
			defer func() {
				observeGetADUsers(time.Since(start), one.Name)
				if r := recover(); r != nil {
					u.log.Error(fmt.Sprintf(" Паника в горутине FillRedisCaсheFromAD для домена %s: %v", one.Name, r))
				}
				wg.Done()
			}()
			u.log.Info("Получение данных домена " + one.Name + "...")
			internetGroups := u.ad.GetDomainInternetGroups(one.Name)

			users, err := u.ad.GetDomainUsers(one.Name)
			comps, _ := u.ad.GetDomainComputers(one.Name)
			rmsPort := u.ad.GetDomainRMSPort(one.Name)
			u.log.Info("Получение данных домена " + one.Name + " завершено. Обработка данных...")
			if err == nil || len(users) > 0 {
				//				println("Get from ad to redis from " + one.Name)
				ips, _ := u.repo.GetDomainUsersIP(one.Name)
				avatars, _ := u.repo.GetDomainUsersAvatars(one.Name)
				for _, user := range users {

					user["domain"] = one
					user["rms_port"] = rmsPort
					if val, ok := user["userAccountControl"]; ok {
						uac := val.(string)
						if uac == "514" || uac == "66050" {
							user["disabled"] = true
							disabledUsersCount++
						} else {
							enabledUsersCount++
						}
					}
					if val, ok := user["msDS-User-Account-Control-Computed"]; ok {
						if val.(string) == "16" {
							user["locked"] = true
							lockedUsersCount++
						}
					}
					if val, ok := user["lockoutTime"]; ok {
						if val.(string) != "0" {
							lockoutTime, err := ADFiletimeToGoTime(val.(string))
							if err == nil {
								user["lockoutTime"] = lockoutTime
							}
						}
					}

					arch := user["distinguishedName"].(string)
					arch = strings.ToUpper(arch)
					if strings.Contains(arch, "=АРХИВ") {
						user["archive"] = true
					}
					userGroups := []string{}
					userGroups, ok := user["memberOf"].([]string)
					if ok {
						userGroups = user["memberOf"].([]string)

					}
					if IsStringInArray("Отключенные Кадровичком", userGroups) {
						user["sap_disabled"] = true
					}

					if AnyOfFirstInSecond(internetGroups.WhiteList, userGroups) {
						user["internetwl"] = true
						whiteListInternetUsersCount++
					}
					if AnyOfFirstInSecond(internetGroups.Full, userGroups) {
						user["internet"] = true
						fullInternetUsersCount++
					}
					if AnyOfFirstInSecond(internetGroups.Tech, userGroups) {
						user["internettech"] = true
						techInternetUsersCount++
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
						u.redis.AddKeyFieldValue("allusers:stagging", upn, jsonUser)
					}

					//				redisClient.HSet(ctx, "allusers", user["userPrincipalName"], jsonUser).Err()
				}

			}
			sort.Slice(users, func(i, j int) bool {
				return fmt.Sprintf("%v", users[i]["cn"]) < fmt.Sprintf("%v", users[j]["cn"])
			})
			observeUsersPerDomain(one.Name, "all", "all", len(users))
			observeUsersPerDomain(one.Name, "all", "enabled", enabledUsersCount)
			observeUsersPerDomain(one.Name, "all", "disabled", disabledUsersCount)
			observeUsersPerDomain(one.Name, "all", "locked", lockedUsersCount)
			observeUsersPerDomain(one.Name, "internet", "full", fullInternetUsersCount)
			observeUsersPerDomain(one.Name, "internet", "white", whiteListInternetUsersCount)
			observeUsersPerDomain(one.Name, "internet", "tech", techInternetUsersCount)
			jsonUsers, _ := json.Marshal(users)
			jsonComps, _ := json.Marshal(comps)
			u.redis.DelKeyField("ad", one.Name)
			err1 := u.redis.AddKeyFieldValue("ad", one.Name, jsonUsers)
			if err1 != nil {
				return
			}
			u.redis.DelKeyField("adc", one.Name)
			u.redis.AddKeyFieldValue("adc", one.Name, jsonComps)
			u.log.Info("Получение данных домена " + one.Name + " завершено. Обработка данных завершена.")
		}(one)
	}
	wg.Wait()
	u.redis.Rename("allusers", "prev")
	u.redis.Rename("allusers:stagging", "allusers")
	u.redis.Unlink("prev")
	u.log.Info("Получение данных всех доменов завершено.")

	return nil
}
