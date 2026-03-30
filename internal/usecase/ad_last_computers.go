package usecase

import (
	"encoding/json"
	"sync"

	"github.com/saygik/go-userinfo/internal/entity"
)

// // GetADLastComputers возвращает список компьютеров домена с последними логинами,
// // обогащённый пользователями и AD‑свойствами из Redis.
// func (u *UseCase) GetADLastComputers(domain string, perms entity.Permissions) ([]entity.DomainComputer, error) {
// 	// Проверяем, что домен существует
// 	if !u.ad.IsDomainExist(domain) {
// 		return nil, u.Error("домен " + domain + " отсутствует в системе")
// 	}

// 	// Проверяем доступ пользователя к домену

// 	accessLevel := u.GetAccessLevelForDomain(&perms, domain)
// 	if accessLevel == "none" || accessLevel == "user" {
// 		return nil, u.Error("у вас нет прав на просмотр информации по домену " + domain)
// 	}

// 	// Получаем список компьютеров из MSSQL
// 	comps, err := u.repo.GetComputerByDomain(domain)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Получаем всех пользователей из Redis (allusers)
// 	oneDomain, err := u.redis.GetKeyFieldValue("ad", domain)
// 	var allUsers []map[string]interface{}
// 	if err := json.Unmarshal([]byte(oneDomain), &allUsers); err != nil {
// 		return nil, fmt.Errorf("ошибка при обновлении кеша, ошибка парсинга JSON: %w", err)
// 	}

// 	// Получаем компьютеры домена из Redis (adc)
// 	domainCompsJSON, err := u.redis.GetKeyFieldValue("adc", domain)
// 	if err != nil {
// 		// Если ключа нет, просто продолжаем без AD‑свойств
// 		domainCompsJSON = ""
// 	}

// 	var redisDomainComps []map[string]interface{}
// 	if domainCompsJSON != "" {
// 		unmarshalString(domainCompsJSON, &redisDomainComps)
// 	}

// 	// Подготавливаем быстрый поиск AD‑компьютеров по имени
// 	adCompByName := map[string]entity.ComputerProperties{}
// 	for _, c := range redisDomainComps {
// 		if name, ok := c["cn"].(string); ok {
// 			properties := entity.ComputerProperties{
// 				OperatingSystem: GetStringFromMap(c, "operatingSystem"),
// 				Description:     GetStringFromMap(c, "description"),
// 			}
// 			adCompByName[name] = properties
// 		}
// 	}

// 	for i, c := range comps {
// 		comps[i].ID = i
// 		// Находим пользователей с совпадающим именем компьютера
// 		usersForComp := []entity.ComputerUser{}
// 		for _, user := range allUsers {

// 			compName, ok := user["computer"].(string)
// 			if !ok || compName != c.Computer {
// 				continue
// 			}
// 			findedUser := entity.ComputerUser{
// 				UPN:         GetStringFromMap(user, "userPrincipalName"),
// 				DisplayName: GetStringFromMap(user, "displayName"),
// 				Company:     GetStringFromMap(user, "company"),
// 				Department:  GetStringFromMap(user, "department"),
// 				Title:       GetStringFromMap(user, "title"),
// 				Mail:        GetStringFromMap(user, "mail"),
// 				Telephone:   GetStringFromMap(user, "telephoneNumber"),
// 				Computer:    c.Computer,
// 				IP:          c.IP,
// 				LastDate:    GetStringFromMap(user, "ip_date"),
// 			}
// 			usersForComp = append(usersForComp, findedUser)
// 		}

// 		if len(usersForComp) > 0 {
// 			comps[i].Users = usersForComp
// 		}

// 		//		Добавляем AD‑свойства компьютера из Redis
// 		if adProps, ok := adCompByName[c.Computer]; ok {
// 			comps[i].OperatingSystem = adProps.OperatingSystem
// 			comps[i].Description = adProps.Description
// 		}

// 	}

// 	return comps, nil
// }

func (u *UseCase) GetADLastComputers(domain string, perms entity.Permissions) ([]entity.DomainComputer, error) {
	// 1. Early returns — минимум проверок
	if !u.ad.IsDomainExist(domain) {
		return nil, u.Error("домен " + domain + " не существует")
	}

	accessLevel := u.GetAccessLevelForDomain(&perms, domain)
	if accessLevel == "none" || accessLevel == "user" {
		return nil, u.Error("нет прав на домен " + domain)
	}

	// 2. 🔥 Параллельная загрузка данных
	var (
		comps            []entity.DomainComputer
		allUsers         []map[string]interface{}
		redisDomainComps []map[string]interface{}
		adCompByName     map[string]entity.ComputerProperties
		wg               sync.WaitGroup
		mu               sync.Mutex
	)

	domainAdminsMap := make(map[string]string)
	localAdminsMap := make(map[string]string)
	computersTickets := make(map[string][]entity.IdName)

	// MSSQL компьютеры
	wg.Add(1)
	go func() {
		defer wg.Done()
		compsMSSQL, err := u.repo.GetComputerByDomain(domain)
		if err == nil {
			mu.Lock()
			comps = compsMSSQL
			mu.Unlock()
		}
	}()

	// 3. 🔥 Redis данные ОДНИМ запросом через HMGet
	wg.Add(1)
	go func() {
		defer wg.Done()
		redisData, err := u.redis.GetKeyFieldsValue("ad", []string{domain})
		if err == nil && len(redisData) > 0 && redisData[0] != nil {
			json.Unmarshal([]byte(redisData[0].(string)), &allUsers)
		}
	}()

	// Redis компьютеры домена
	wg.Add(1)
	go func() {
		defer wg.Done()
		redisData, err := u.redis.GetKeyFieldsValue("adc", []string{domain})
		if err == nil && len(redisData) > 0 && redisData[0] != nil {
			unmarshalString(redisData[0].(string), &redisDomainComps)
		}
	}()
	if isAdminDomain(perms, domain) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			compsLocalAdmins, _ := u.ComputerLocalAdminsGet(false) // domain = false
			for _, admin := range compsLocalAdmins {
				localAdminsMap[admin.Computer] = admin.Administrators
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			compsLocalAdminsDomain, _ := u.ComputerLocalAdminsGet(true) // domain = true
			for _, admin := range compsLocalAdminsDomain {
				domainAdminsMap[admin.Computer] = admin.Administrators
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			computersTickets, _ = u.glpi.GetComputersTickets()

		}()
	}
	wg.Wait()

	if len(comps) == 0 {
		return comps, nil
	}

	// 4. 🔥 Быстрый поисковый индекс (O(1))
	adCompByName = make(map[string]entity.ComputerProperties, len(redisDomainComps))
	userByComputer := make(map[string][]entity.ComputerUser, len(allUsers))

	// Индексация AD компьютеров O(n)
	for _, c := range redisDomainComps {
		if name, ok := c["cn"].(string); ok {
			adCompByName[name] = entity.ComputerProperties{
				OperatingSystem: GetStringFromMap(c, "operatingSystem"),
				Description:     GetStringFromMap(c, "description"),
			}
		}
	}

	// 🔥 ИНДЕКСАЦИЯ ПОЛЬЗОВАТЕЛЕЙ O(n) вместо O(n²)!
	for _, user := range allUsers {
		compName, ok := user["computer"].(string)
		if !ok || compName == "" {
			continue
		}

		compUser := entity.ComputerUser{
			UPN:         GetStringFromMap(user, "userPrincipalName"),
			DisplayName: GetStringFromMap(user, "displayName"),
			Company:     GetStringFromMap(user, "company"),
			Department:  GetStringFromMap(user, "department"),
			Title:       GetStringFromMap(user, "title"),
			Mail:        GetStringFromMap(user, "mail"),
			Telephone:   GetStringFromMap(user, "telephoneNumber"),
			Computer:    compName,
			IP:          GetStringFromMap(user, "ip"),
			LastDate:    GetStringFromMap(user, "ip_date"),
		}

		userByComputer[compName] = append(userByComputer[compName], compUser)
	}

	// 5. 🔥 Привязка O(n) вместо O(n²)
	for i := range comps {
		comps[i].ID = i
		compName := comps[i].Computer
		// 🔥 O(1) поиск пользователей по имени компьютера
		if users, ok := userByComputer[compName]; ok {
			comps[i].Users = users

			// Заполняем IP и дату из первого пользователя
			if len(users) > 0 {
				comps[i].IP = users[0].IP
			}
		}
		if isAdminDomain(perms, domain) {
			if admins, ok := domainAdminsMap[compName]; ok {
				comps[i].AdministratorsDomain = admins
			}
			// Локальные админы
			if admins, ok := localAdminsMap[compName]; ok {
				comps[i].AdministratorsLocal = admins
			}
		}
		// 🔥 O(1) AD свойства
		if props, ok := adCompByName[compName]; ok {
			comps[i].OperatingSystem = props.OperatingSystem
			comps[i].Description = props.Description
		}
		if ts, ok := computersTickets[compName]; ok {
			comps[i].Tickets = ts
		}
	}

	return comps, nil
}

func isAdminDomain(perms entity.Permissions, domain string) bool {
	return perms.IsAdmin || perms.IsSysAdmin || perms.AdminDomains[domain]
}
