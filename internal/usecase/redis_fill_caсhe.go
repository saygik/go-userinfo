package usecase

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/saygik/go-userinfo/internal/entity"
	"github.com/saygik/go-userinfo/internal/state"
)

// соответствие build-номера Windows маркетинговому имени версии
var windowsBuildToHuman = map[int]string{
	26200: "25H2", // Windows 11 24H2
	26100: "24H2", // Windows 11 24H2
	22631: "23H2", // Windows 11 23H2
	22621: "22H2", // Windows 11 22H2
	22000: "21H2", // Windows 11 21H2

	19045: "22H2", // Windows 10 22H2
	19044: "21H2", // Windows 10 21H2
	19043: "21H1", // Windows 10 21H1
	19042: "20H2", // Windows 10 20H2
	19041: "2004",
	18363: "1909",
	18362: "1903",
	17763: "1809",
	17134: "1803",
	16299: "1709",
	15063: "1703",
	14393: "1607",
	10586: "1511",
	10240: "1507",
	7601:  "WIN7", // Windows 7 SP1
}

// windowsVersionToHuman преобразует строку operatingSystemVersion AD в вид 24H2 / 22H2 и т.п.
// Примеры входа:
//   - "10.0 (26100)" -> "24H2"
//   - "10.0 (19045)" -> "22H2"
func windowsVersionToHuman(osVersion string) string {
	if osVersion == "" {
		return ""
	}

	// оставляем только числовые последовательности и берём последнюю (build)
	parts := strings.FieldsFunc(osVersion, func(r rune) bool {
		return r < '0' || r > '9'
	})
	if len(parts) == 0 {
		return ""
	}

	buildStr := parts[len(parts)-1]
	build, err := strconv.Atoi(buildStr)
	if err != nil {
		return ""
	}

	if human, ok := windowsBuildToHuman[build]; ok {
		return human
	}

	return ""
}

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

	// Копируем данные из текущего allusers в stagging, чтобы сохранить данные для доменов, которые не успеют обновиться
	u.redis.Delete("allusers:stagging")
	existingUsers, err := u.redis.GetKeyFieldAll("allusers")
	if err == nil && len(existingUsers) > 0 {
		// Копируем все существующие данные в stagging
		for upn, userData := range existingUsers {
			u.redis.AddKeyFieldValue("allusers:stagging", upn, []byte(userData))
		}
		u.log.Info(fmt.Sprintf("Скопировано %d пользователей из allusers в allusers:stagging для сохранения данных при возможном таймауте", len(existingUsers)))
	}

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
			// Добавляем человекочитаемую версию ОС по operatingSystemVersion (например, 24H2)
			for _, comp := range comps {
				if v, ok := comp["operatingSystemVersion"]; ok {
					if verStr, ok := v.(string); ok {
						if human := windowsVersionToHuman(verStr); human != "" {
							comp["operatingSystemVersionHuman"] = human
						}
					}
				}
			}
			rmsPort := u.ad.GetDomainRMSPort(one.Name)
			u.log.Info("Получение данных домена " + one.Name + " завершено. Обработка данных...")
			if err == nil || len(users) > 0 {
				// Удаляем старых пользователей этого домена из stagging перед добавлением новых
				// Это нужно, чтобы удаленные пользователи не оставались в кеше
				// Делаем это только после успешного получения данных, чтобы при ошибке сохранить старые данные
				stagingUsers, _ := u.redis.GetKeyFieldAll("allusers:stagging")
				deletedCount := 0
				for upn, userDataStr := range stagingUsers {
					var userData map[string]interface{}
					if err := json.Unmarshal([]byte(userDataStr), &userData); err == nil {
						if domainObj, ok := userData["domain"].(map[string]interface{}); ok {
							if domainName, ok := domainObj["name"].(string); ok && domainName == one.Name {
								u.redis.DelKeyField("allusers:stagging", upn)
								deletedCount++
							}
						}
					}
				}
				if deletedCount > 0 {
					u.log.Info(fmt.Sprintf("Удалено %d старых пользователей домена %s из stagging", deletedCount, one.Name))
				}
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

			if len(users) > 0 {
				jsonUsers, _ := json.Marshal(users)
				u.redis.DelKeyField("ad", one.Name)
				err1 := u.redis.AddKeyFieldValue("ad", one.Name, jsonUsers)
				if err1 != nil {
					return
				}
				u.log.Info(fmt.Sprintf("Добавлено %d пользователей домена %s", len(users), one.Name))
			}
			if len(comps) > 0 {
				jsonComps, _ := json.Marshal(comps)
				u.redis.DelKeyField("adc", one.Name)
				u.redis.AddKeyFieldValue("adc", one.Name, jsonComps)
				u.log.Info(fmt.Sprintf("Добавлено %d компьютеров домена %s", len(comps), one.Name))

			}
			u.log.Info(fmt.Sprintf("Получение данных домена  %s завершено. Обработка данных завершена.", one.Name))

		}(one)
	}

	// Ожидание завершения всех горутин с таймаутом
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Таймаут ожидания - 8 минут для всех доменов (меньше интервала запуска 10 минут)
	timeout := 2 * time.Minute
	select {
	case <-done:
		u.log.Info("Все горутины доменов завершены успешно.")
	case <-time.After(timeout):
		u.log.Error(fmt.Sprintf("Таймаут ожидания завершения обработки доменов (превышен лимит %v). Продолжаем выполнение.", timeout))
	}

	u.redis.Rename("allusers", "prev")
	u.redis.Rename("allusers:stagging", "allusers")
	u.redis.Unlink("prev")
	u.log.Info("Получение данных всех доменов завершено.")

	return nil
}
