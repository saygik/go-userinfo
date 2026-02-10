package usecase

import (
	"sort"

	"github.com/saygik/go-userinfo/internal/entity"
)

// GetADComputersOSFamily возвращает количество компьютеров домена, сгруппированное по семейству ОС (OperatingSystemFamily).
func (u *UseCase) GetADComputersOSFamily(domain, user string) ([]entity.ComputerFamilyCount, error) {

	access := u.GetAccessToResource(domain, user)
	if access == -1 {
		return nil, u.Error("у вас нет прав на просмотр информации по домену " + domain)
	}

	if !u.ad.IsDomainExist(domain) {
		return nil, u.Error("домен " + domain + " отсутствует в системе")
	}

	if domain == "все домены" {

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

	return res, nil
}
