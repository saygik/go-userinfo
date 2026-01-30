package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/saygik/go-userinfo/internal/entity"
)

// GetADLastComputers возвращает список компьютеров домена с последними логинами,
// обогащённый пользователями и AD‑свойствами из Redis.
func (u *UseCase) GetADLastComputers(domain, user string) ([]entity.DomainComputer, error) {
	// Проверяем, что домен существует
	if !u.ad.IsDomainExist(domain) {
		return nil, u.Error("домен " + domain + " отсутствует в системе")
	}

	// Проверяем доступ пользователя к домену
	access := u.GetAccessToResource(domain, user)
	if access == -1 {
		return nil, u.Error("у вас нет прав на просмотр информации по домену " + domain)
	}

	// Получаем список компьютеров из MSSQL
	comps, err := u.repo.GetComputerByDomain(domain)
	if err != nil {
		return nil, err
	}

	// Получаем всех пользователей из Redis (allusers)
	oneDomain, err := u.redis.GetKeyFieldValue("ad", domain)
	var allUsers []map[string]interface{}
	if err := json.Unmarshal([]byte(oneDomain), &allUsers); err != nil {
		return nil, fmt.Errorf("ошибка при обновлении кеша, ошибка парсинга JSON: %w", err)
	}

	// Получаем компьютеры домена из Redis (adc)
	domainCompsJSON, err := u.redis.GetKeyFieldValue("adc", domain)
	if err != nil {
		// Если ключа нет, просто продолжаем без AD‑свойств
		domainCompsJSON = ""
	}

	var redisDomainComps []map[string]interface{}
	if domainCompsJSON != "" {
		unmarshalString(domainCompsJSON, &redisDomainComps)
	}

	// Подготавливаем быстрый поиск AD‑компьютеров по имени
	adCompByName := map[string]entity.ComputerProperties{}
	for _, c := range redisDomainComps {
		if name, ok := c["cn"].(string); ok {
			properties := entity.ComputerProperties{
				OperatingSystem: GetStringFromMap(c, "operatingSystem"),
				Description:     GetStringFromMap(c, "description"),
			}
			adCompByName[name] = properties
		}
	}

	for i, c := range comps {
		comps[i].ID = i
		// Находим пользователей с совпадающим именем компьютера
		usersForComp := []entity.ComputerUser{}
		for _, user := range allUsers {

			compName, ok := user["computer"].(string)
			if !ok || compName != c.Computer {
				continue
			}
			findedUser := entity.ComputerUser{
				UPN:         GetStringFromMap(user, "userPrincipalName"),
				DisplayName: GetStringFromMap(user, "displayName"),
				Company:     GetStringFromMap(user, "company"),
				Department:  GetStringFromMap(user, "department"),
				Title:       GetStringFromMap(user, "title"),
				Mail:        GetStringFromMap(user, "mail"),
				Telephone:   GetStringFromMap(user, "telephoneNumber"),
				Computer:    c.Computer,
				IP:          c.IP,
				LastDate:    GetStringFromMap(user, "ip_date"),
			}
			usersForComp = append(usersForComp, findedUser)
		}

		if len(usersForComp) > 0 {
			comps[i].Users = usersForComp
		}

		//		Добавляем AD‑свойства компьютера из Redis
		if adProps, ok := adCompByName[c.Computer]; ok {
			comps[i].OperatingSystem = adProps.OperatingSystem
			comps[i].Description = adProps.Description
		}

	}

	return comps, nil
}
