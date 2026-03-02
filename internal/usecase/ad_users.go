package usecase

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/saygik/go-userinfo/internal/entity"
)

var sensitiveFields = []string{
	"ip", "pwdLastSet", "proxyAddresses", "passwordDontExpire",
	"passwordCantChange", "distinguishedName", "userAccountControl", "memberOf",
}

var sensitiveLongFields = []string{
	"ip", "pwdLastSet", "proxyAddresses", "passwordDontExpire",
	"passwordCantChange", "distinguishedName", "userAccountControl", "memberOf", "presence", "url", "otherTelephone",
}

// func (u *UseCase) GetADUsers(perms entity.Permissions) ([]map[string]interface{}, error) {

// 	redisADUsers, err := u.redis.GetKeyFieldAll("ad")
// 	if err != nil {
// 		return nil, err
// 	}

// 	var res []map[string]interface{}
// 	for domainName, oneDomain := range redisADUsers {

// 		accessLevel := u.GetAccessLevelForDomain(&perms, domainName)
// 		if accessLevel == "none" {
// 			continue
// 		}

// 		var r []map[string]interface{}
// 		json.Unmarshal([]byte(oneDomain), &r)

// 		for _, user := range r {

// 			if accessLevel == "user" {
// 				delete(user, "ip")
// 				delete(user, "pwdLastSet")
// 				delete(user, "proxyAddresses")
// 				delete(user, "passwordDontExpire")
// 				delete(user, "passwordCantChange")
// 				delete(user, "distinguishedName")
// 				delete(user, "userAccountControl")
// 				delete(user, "memberOf")

// 				user["restricted"] = true
// 			}
// 			// if !accessToTechnicalInfo && accessToShortTechnicalInfo {
// 			// 	delete(user, "pwdLastSet")
// 			// 	delete(user, "proxyAddresses")
// 			// 	delete(user, "passwordDontExpire")
// 			// 	delete(user, "passwordCantChange")
// 			// 	delete(user, "distinguishedName")
// 			// 	delete(user, "userAccountControl")
// 			// 	delete(user, "memberOf")

// 			// 	user["partial_restricted"] = true
// 			// }
// 		}
// 		res = append(res, r...)
// 	}
// 	return res, nil
// }

func (u *UseCase) GetADUsers(perms entity.Permissions) ([]map[string]interface{}, error) {
	// 1. 🔥 КЭШ с учетом прав (5 минут)
	cacheKey := u.getUsersCacheKey(&perms)

	if cached, err := u.redis.GetKeyValue(cacheKey); err == nil {
		var res []map[string]interface{}
		if json.Unmarshal([]byte(cached), &res) == nil {
			return res, nil
		}
	}

	// 2. 🔥 Только разрешенные домены (HMGet)
	allowedDomains := u.getAllowedUserDomains(&perms)
	if len(allowedDomains) == 0 {
		return []map[string]interface{}{}, nil
	}

	// 3. 🔥 1 Redis запрос!
	redisData, err := u.redis.GetKeyFieldsValue("ad", allowedDomains)
	if err != nil {
		return nil, fmt.Errorf("redis users: %w", err)
	}

	// 4. 🔥 Параллельная обработка доменов
	var result []map[string]interface{}
	var wg sync.WaitGroup
	mu := sync.RWMutex{}

	for i, data := range redisData {
		if data == nil {
			continue
		}

		domain := allowedDomains[i]
		accessLevel := u.GetAccessLevelForDomain(&perms, domain)
		wg.Add(1)

		go func(domain string, dataStr string, level string) {
			defer wg.Done()

			users, err := u.processDomainUsers([]byte(dataStr), level)
			if err == nil && len(users) > 0 {
				mu.Lock()
				result = append(result, users...)
				mu.Unlock()
			}
		}(domain, data.(string), accessLevel)
	}

	wg.Wait()

	// 5. 🔥 Асинхронный кэш
	go u.cacheUsersResult(cacheKey, result)

	return result, nil
}

// 🔥 Функция обработки пользователей домена
func (u *UseCase) processDomainUsers(data []byte, accessLevel string) ([]map[string]interface{}, error) {
	var users []map[string]interface{}
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	// 🔥 Преаллокация слайса
	result := make([]map[string]interface{}, 0, len(users))

	for _, user := range users {
		processedUser := u.cloneUserMap(user) // Копируем чтобы не мутировать оригинал

		if accessLevel == "user" {
			for _, field := range sensitiveFields {
				delete(processedUser, field)
			}
			processedUser["restricted"] = true
		}

		result = append(result, processedUser)
	}

	return result, nil
}

// 🔥 Быстрое копирование map
func (u *UseCase) cloneUserMap(user map[string]interface{}) map[string]interface{} {
	cloned := make(map[string]interface{}, len(user))
	for k, v := range user {
		cloned[k] = v
	}
	return cloned
}

// 🔥 Ключ кэша по правам
func (u *UseCase) getUsersCacheKey(perms *entity.Permissions) string {
	// Сериализуем права для уникального ключа
	keyParts := []string{"users"}
	for _, role := range perms.Roles {
		keyParts = append(keyParts, role.Name)
	}
	for domain := range perms.AdminDomains {
		keyParts = append(keyParts, "A:"+domain)
	}
	for domain := range perms.TechDomains {
		keyParts = append(keyParts, "T:"+domain)
	}
	return strings.Join(keyParts, ":")
}

// 🔥 Разрешенные домены для пользователей
func (u *UseCase) getAllowedUserDomains(perms *entity.Permissions) []string {
	var domains []string

	domainsList := u.ad.DomainList()
	for _, oneDomain := range domainsList {
		accessLevel := u.GetAccessLevelForDomain(perms, oneDomain.Name)
		if accessLevel != "none" {
			domains = append(domains, oneDomain.Name)
		}
	}
	return domains
}

// 🔥 Кэширование результата
func (u *UseCase) cacheUsersResult(cacheKey string, users []map[string]interface{}) {
	data, _ := json.Marshal(users)
	u.redis.AddKeyValue(cacheKey, data, 5*time.Minute)
}
