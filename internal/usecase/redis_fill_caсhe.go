package usecase

import (
	"context"
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

func (u *UseCase) FillRedisCaсheFromAD(ctx context.Context) error {
	if state.IsFillingRedis() {
		u.log.Error("кеш не обновлён, попытка заполнения кеша пользователей при невыполненной предыдущей задаче...")
		return u.Error("кеш не обновлён")
	}
	defer func() {
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

	// Получаем локальных и доменных админов
	compsLocalAdminsDomain, _ := u.ComputerLocalAdminsGet(true) // domain = true
	compsLocalAdmins, _ := u.ComputerLocalAdminsGet(false)      // domain = false

	// Создаем карту для быстрого поиска
	domainAdminsMap := make(map[string]string)
	localAdminsMap := make(map[string]string)

	for _, admin := range compsLocalAdminsDomain {
		domainAdminsMap[admin.Computer] = admin.Administrators
	}
	for _, admin := range compsLocalAdmins {
		localAdminsMap[admin.Computer] = admin.Administrators
	}

	var wg sync.WaitGroup
	for _, one := range adl {
		wg.Add(1)
		go func(one entity.DomainList) {

			enabledUsersCount, disabledUsersCount, lockedUsersCount := 0, 0, 0
			fullInternetUsersCount, whiteListInternetUsersCount, techInternetUsersCount := 0, 0, 0
			serviceUsersCount, serviceMailUsersCount, privilegedUsersCount := 0, 0, 0
			archiveUsersCount, sapDisabledUsersCount, attributeErrorsUsersCount := 0, 0, 0
			passwordDontExpireUsersCount, passwordCantChangeUsersCount := 0, 0
			noComputerUsersCount, osUnsupportedUsersCount := 0, 0
			result := "success"
			start := time.Now()
			defer func() {
				if r := recover(); r != nil {
					result = "panic"
					observeADCacheFillError(one.Name, "panic")
					u.log.Error(fmt.Sprintf(" Паника в горутине FillRedisCaсheFromAD для домена %s: %v", one.Name, r))
				}
				d := time.Since(start)
				observeGetADUsers(d, one.Name)
				observeADCacheFillDuration(d, one.Name, result)
				wg.Done()
			}()
			u.log.Info("Получение данных домена " + one.Name + "...")
			internetGroups := u.ad.GetDomainInternetGroups(one.Name)
			serviceGroups := u.ad.GetDomainServiceGroups(one.Name)

			users, err := u.ad.GetDomainUsers(one.Name)
			if err != nil {
				result = "error"
				if len(users) > 0 {
					result = "partial"
				}
				observeADCacheFillError(one.Name, "ad")
			}
			comps, _ := u.ad.GetDomainComputers(one.Name)
			tagsMap, _ := u.glpi.GetComputersTags(one.Name)
			arhsMap, _ := u.glpi.GetComputersArhs(one.Name)
			computersOs := make(map[string]string)
			computersOsVersion := make(map[string]WindowsVersionInfo)

			// Добавляем человекочитаемую версию ОС по operatingSystemVersion (например, 24H2)
			for _, comp := range comps {
				compName := comp["name"].(string)
				// Обработка dNSHostName
				if dnsHost, ok := comp["dNSHostName"].(string); ok && dnsHost != "" {
					// Изменяем значение dNSHostName: оставляем только часть до точки
					if before, _, ok := strings.Cut(dnsHost, "."); ok {
						comp["dNSHostName"] = before
					}

					// Проверяем совпадает ли dNSHostName (после обрезки) с compName
					comp["longName"] = !strings.EqualFold(comp["dNSHostName"].(string), compName)
				}
				// Доменные админы
				if admins, ok := domainAdminsMap[compName]; ok {
					comp["administrators_domain"] = admins
				}
				// Локальные админы
				if admins, ok := localAdminsMap[compName]; ok {
					comp["administrators_local"] = admins
				}
				//											computersOs[strings.ToUpper(compName)] = comp["operatingSystem"].(string)
				if v, ok := comp["operatingSystem"]; ok {
					if verStr, ok := v.(string); ok {
						computersOs[strings.ToUpper(compName)] = verStr
					}
				}
				if v, ok := comp["operatingSystemVersion"]; ok {
					if verStr, ok := v.(string); ok {
						osName, _ := comp["operatingSystem"].(string)
						if human := windowsVersionToHumanWithSupport(osName, verStr); human.Version != "" {
							comp["operatingSystemVersionHuman"] = human.Version
							computersOsVersion[strings.ToUpper(compName)] = human
						}
					}
				}
				if arh, ok := arhsMap[strings.ToUpper(compName)]; ok {
					comp["computer_arh"] = arh
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

					userMail, ok := user["mail"].(string)
					if ok {
						if strings.ContainsAny(userMail, " \t\n\r") {
							// Неверный email - содержит пробелы внутри
							user["errors"] = "Ошибка в email: содержит пробелы или табуляцию"
						}
					}

					userGroups := []string{}
					userGroups, ok = user["memberOf"].([]string)
					if ok {
						userGroups = user["memberOf"].([]string)
					}
					if IsStringInArray("Отключенные Кадровичком", userGroups) {
						user["sap_disabled"] = true
					}

					if isServiceAccountUser(userGroups, serviceGroups) {
						user["service"] = true
					}

					if isServiceMailUser(userGroups) {
						user["serviceMail"] = true
					}

					if isPrivilegedUser(userGroups) {
						user["privileged"] = true
					}

					if user["service"] == true {
						serviceUsersCount++
					}
					if user["serviceMail"] == true {
						serviceMailUsersCount++
					}
					if user["privileged"] == true {
						privilegedUsersCount++
					}
					if user["archive"] == true {
						archiveUsersCount++
					}
					if user["sap_disabled"] == true {
						sapDisabledUsersCount++
					}
					if userErrors, ok := user["errors"].(string); ok && userErrors != "" {
						attributeErrorsUsersCount++
					}
					if val, ok := user["userAccountControl"]; ok {
						if uacStr, ok := val.(string); ok {
							if uac, err := strconv.Atoi(uacStr); err == nil {
								if uac&65536 == 65536 {
									passwordDontExpireUsersCount++
								}
								if uac&64 == 64 {
									passwordCantChangeUsersCount++
								}
							}
						}
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
								comp := strings.ToUpper(ip.Computer)
								user["ip"] = ip.Ip
								user["computer"] = ip.Computer
								user["ip_date"] = ip.IpDate
								user["rms_installed"] = ip.Rms
								if tags, ok := tagsMap[comp]; ok {
									user["computer_glpi_tags"] = tags
								}
								if tags, ok := arhsMap[comp]; ok {
									user["computer_arh"] = tags
								}
								if os, ok := computersOs[comp]; ok {
									user["computer_os"] = os
								}
								if cosv, ok := computersOsVersion[comp]; ok {
									user["computer_os_version"] = cosv
								}
							}
						}
					}
					if _, ok := user["computer"].(string); !ok {
						noComputerUsersCount++
					} else if version, ok := user["computer_os_version"].(WindowsVersionInfo); !ok || !version.Supported {
						osUnsupportedUsersCount++
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
			observeUsersPerDomain(one.Name, "privileged", "all", privilegedUsersCount)
			observeUsersPerDomain(one.Name, "service_mail", "all", serviceMailUsersCount)
			observeUsersPerDomain(one.Name, "service", "all", serviceUsersCount)
			observeUsersPerDomain(one.Name, "archive", "all", archiveUsersCount)
			observeUsersPerDomain(one.Name, "sap_disabled", "all", sapDisabledUsersCount)
			observeUsersPerDomain(one.Name, "attribute_errors", "all", attributeErrorsUsersCount)
			observeUsersPerDomain(one.Name, "password_never_expires", "all", passwordDontExpireUsersCount)
			observeUsersPerDomain(one.Name, "password_cant_change", "all", passwordCantChangeUsersCount)
			observeUsersPerDomain(one.Name, "no_computer", "all", noComputerUsersCount)
			observeUsersPerDomain(one.Name, "os_unsupported", "all", osUnsupportedUsersCount)

			if len(users) > 0 {
				jsonUsers, _ := json.Marshal(users)
				u.redis.DelKeyField("ad", one.Name)
				err1 := u.redis.AddKeyFieldValue("ad", one.Name, jsonUsers)
				if err1 != nil {
					result = "error"
					observeADCacheFillError(one.Name, "redis")
					return
				}
				observeADCacheLastSuccess(one.Name)
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
		observeADCacheLastComplete()
	case <-ctx.Done():
		u.log.Info("Заполнение кеша AD прервано: сервис останавливается.")
		return nil
	case <-time.After(timeout):
		observeADCacheFillTimeout()
		u.log.Error(fmt.Sprintf("Таймаут ожидания завершения обработки доменов (превышен лимит %v). Продолжаем выполнение.", timeout))
	}

	if !state.IsInitialized() {
		state.SetInitialized()
	}

	u.redis.Rename("allusers", "prev")
	u.redis.Rename("allusers:stagging", "allusers")
	u.redis.Unlink("prev")
	state.SetLastADFill(time.Now())
	u.log.Info("Получение данных всех доменов завершено.")

	return nil
}
