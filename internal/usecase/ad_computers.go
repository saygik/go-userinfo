package usecase

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/saygik/go-userinfo/internal/entity"
)

// func (u *UseCase) GetADComputers(perms entity.Permissions) ([]map[string]interface{}, error) {

// 	var res []map[string]interface{}
// 	redisADComputers, err := u.redis.GetKeyFieldAll("adc")
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, oneDomain := range redisADComputers {

// 		accessLevel := u.GetAccessLevelForDomain(&perms, oneDomain)
// 		if accessLevel != "none" || accessLevel != "user" {
// 			var r []map[string]interface{}
// 			unmarshalString(oneDomain, &r)

// 			res = append(res, r...)
// 		}

// 	}
// 	return res, nil
// }

func (u *UseCase) GetADComputers(perms entity.Permissions) ([]map[string]interface{}, error) {

	// comp["administrators_domain"] = admins
	// comp["administrators_local"] = admins
	allowedDomains := u.getAllowedDomainsForComputers(&perms)
	if len(allowedDomains) == 0 {
		return nil, nil
	}

	redisData, err := u.redis.GetKeyFieldsValue("adc", allowedDomains)
	if err != nil {
		return nil, fmt.Errorf("redis HMGet: %w", err)
	}
	isAdmin := perms.IsAdmin || perms.IsSysAdmin
	// Параллельный парсинг
	var result []map[string]any
	var wg sync.WaitGroup
	mu := sync.Mutex{}

	for _, data := range redisData {
		if data == nil {
			continue
		}
		//		domain := allowedDomains[i]
		wg.Add(1)
		go func(dataStr string) {
			defer wg.Done()
			var computers []map[string]any
			if json.Unmarshal([]byte(dataStr), &computers) == nil {
				for i := range computers {
					comp := &computers[i] // ссылка для изменения
					if isAdmin {
						continue
					}
					if domain, ok := (*comp)["domain"].(string); ok && perms.AdminDomains[domain] {
						continue
					}
					// Удаляем поля для не-админов
					delete(*comp, "administrators_domain")
					delete(*comp, "administrators_local")
					delete(*comp, "servicePrincipalName")
					delete(*comp, "ms-Mcs-AdmPwdExpirationTime")

				}
				mu.Lock()
				result = append(result, computers...)
				mu.Unlock()
			}
		}(data.(string))
	}

	wg.Wait()
	return result, nil
}

func (u *UseCase) getAllowedDomainsForComputers(perms *entity.Permissions) []string {
	var domains []string

	domainsList := u.ad.DomainList()
	for _, oneDomain := range domainsList {
		accessLevel := u.GetAccessLevelForDomain(perms, oneDomain.Name)
		if accessLevel != "none" && accessLevel != "user" {
			domains = append(domains, oneDomain.Name)
		}
	}
	return domains // Только admin/tech для компьютеров!
}
