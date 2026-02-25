package usecase

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/saygik/go-userinfo/internal/entity"
)

// GetADComputersOSFamily возвращает количество компьютеров домена, сгруппированное по семейству ОС (OperatingSystemFamily).
func (u *UseCase) GetADComputersOSFamily(domain string, perms entity.Permissions) ([]entity.ComputerFamilyCount, error) {

	accessLevel := u.GetAccessLevelForDomain(&perms, domain)
	if accessLevel == "none" || accessLevel == "user" {
		return nil, u.Error("нет прав на домен " + domain)
	}

	if !u.ad.IsDomainExist(domain) {
		return nil, u.Error("домен " + domain + " отсутствует в системе")
	}
	// 1. КЭШ Redis (5 минут)
	cacheKey := fmt.Sprintf("os_family:%s", domain)
	if cached, err := u.redis.GetKeyValue(cacheKey); err == nil {
		var res []entity.ComputerFamilyCount
		if json.Unmarshal([]byte(cached), &res) == nil && len(res) > 0 {
			return res, nil
		}
	}

	comps, err := u.ad.GetDomainComputers(domain)
	if err != nil {
		return nil, err
	}

	counts := map[string]int{}
	for _, c := range comps {
		osName, _ := c["operatingSystem"].(string)
		family := osFamilyName(osName)
		if family == "" {
			// Неопознанные/нестандартные ОС пропускаем
			continue
		}
		counts[family]++
	}

	res := make([]entity.ComputerFamilyCount, 0, len(counts))
	for fam, cnt := range counts {
		res = append(res, entity.ComputerFamilyCount{
			OperatingSystemFamily: fam,
			Count:                 cnt,
		})
	}

	// Сортируем по тому же порядку семейств, что и версии
	sort.Slice(res, func(i, j int) bool {
		fi := osFamilyRank(res[i].OperatingSystemFamily)
		fj := osFamilyRank(res[j].OperatingSystemFamily)
		if fi != fj {
			return fi < fj
		}
		// при равном ранге — по названию
		return res[i].OperatingSystemFamily < res[j].OperatingSystemFamily
	})
	// 6. 🔥 Кэш результата асинхронно
	go func() {
		data, _ := json.Marshal(res)
		u.redis.AddKeyValue(cacheKey, data, 30*time.Minute)
	}()

	return res, nil
}
