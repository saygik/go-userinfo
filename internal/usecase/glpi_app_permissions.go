package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) LoadUserPermissions(upn string) (entity.Permissions, error) {
	perms := entity.Permissions{
		AdminDomains: make(map[string]bool),
		TechDomains:  make(map[string]bool),
		UserDomains:  make(map[string]bool),
		User:         upn,
	}
	// 1. Все домены системы
	domains := u.ad.DomainList()
	// 2. Роли пользователя
	roles, err := u.repo.GetUserRoles(upn)
	if err == nil {
		perms.Roles = roles
		for _, role := range roles {
			switch role.Name {
			case "DomainUser":
				perms.AllDomains = true
			case "SysAdmin":
				perms.IsSysAdmin = true
				perms.IsAdmin = true
				perms.IsTech = true
			case "Admin":
				perms.IsAdmin = true
				perms.IsTech = true
			case "Tech":
				perms.IsTech = true
			}
		}
	}

	// 3. Домашний домен из UPN
	perms.HomeDomain = getDomainFromUserName(upn)

	// 4. user_domain_access права
	domainAccess, err := u.repo.GetDomainAccess(upn)
	if err == nil {
		for _, da := range domainAccess {
			switch da.AccessLevel {
			case "admin":
				perms.AdminDomains[da.Domain] = true
			case "tech":
				perms.TechDomains[da.Domain] = true
			case "user":
				perms.UserDomains[da.Domain] = true
			}
		}
	}
	perms.DomainsPermission = u.buildDomainsList(&perms, domains)
	// 5. Дополнить home_domain
	if perms.HomeDomain != "" {
		perms.UserDomains[perms.HomeDomain] = true
		if perms.IsTech {
			perms.TechDomains[perms.HomeDomain] = true
		}
	}

	// 6. Расширить права по ролям для ВСЕХ доменов
	for _, domain := range domains {
		domainName := domain.Name

		if perms.AllDomains {
			perms.UserDomains[domainName] = true
		}
		if perms.IsAdmin {
			perms.TechDomains[domainName] = true
		}
		if perms.IsSysAdmin {
			perms.AdminDomains[domainName] = true
		}
	}
	perms.Domains = u.buildDomainsList(&perms, domains)
	// 7. Разделы сайта из user_permissions
	sections, err := u.repo.GetSections(upn)
	if err == nil {
		perms.Sections = sections
	}

	// 8. Авто-добавление AD разделов
	if len(perms.AdminDomains) > 0 || len(perms.TechDomains) > 0 || len(perms.UserDomains) > 0 {
		perms.Sections = appendUniqueIdName(perms.Sections, entity.IdNameDescription{Id: 1, Name: "/ad/users", Description: "Пользователи AD"})

	}
	if len(perms.TechDomains) > 0 || len(perms.AdminDomains) > 0 {
		adUsersSection := []entity.IdNameDescription{
			{Id: 2, Name: "/ad/computers", Description: "Компьютеры AD"},
			{Id: 3, Name: "/ad/lastcomputers", Description: "Активные компьютеры AD"},
			{Id: 4, Name: "/ad/stat", Description: "Статистика AD"},
		}
		perms.Sections = appendUniqueIdName(perms.Sections, adUsersSection...)
	}

	return perms, nil
}
func appendUniqueIdName(slice []entity.IdNameDescription, items ...entity.IdNameDescription) []entity.IdNameDescription {
	seen := make(map[string]bool)

	// Существующие
	for _, item := range slice {
		seen[item.Name] = true
	}

	// Новые без дублей
	for _, item := range items { // 🔥 items это variadic (0+ элементов)
		if !seen[item.Name] {
			slice = append(slice, item)
			seen[item.Name] = true
		}
	}

	return slice
}

func (u *UseCase) buildDomainsList(perms *entity.Permissions, allDomains []entity.DomainList) []entity.IdNameDescription {
	result := make([]entity.IdNameDescription, 0, len(allDomains))
	domainCounter := 1 // ID начинается с 1

	for _, domain := range allDomains {
		maxAccessLevel := u.GetAccessLevelForDomain(perms, domain.Name)

		// 🔥 Максимальный уровень доступа для Description
		description := maxAccessLevel

		result = append(result, entity.IdNameDescription{
			Id:          domainCounter,
			Name:        domain.Name,
			Description: description,
		})
		domainCounter++
	}

	return result
}

// func (u *UseCase) LoadUserPermissions(upn string) (entity.Permissions, error) {
// 	perms := entity.Permissions{}

// 	// Параллельная загрузка через goroutines (опционально)
// 	var wg sync.WaitGroup
// 	wg.Add(4)

// 	go func() {
// 		defer wg.Done()
// 		roles, err := u.repo.GetRoles(upn)
// 		if err == nil {
// 			perms.Roles = roles
// 			for _, role := range roles {
// 				if role == "DomainUser" {
// 					perms.AllDomains = true
// 				}
// 			}
// 		}
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		domains, err := u.repo.GetUserDomains(upn)
// 		if err == nil {
// 			perms.TechDomains = domains
// 		}
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		sections, err := u.repo.GetUserSections(upn)
// 		if err == nil {
// 			perms.Sections = sections
// 		}
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		homeDomain := getDomainFromUserName(upn)
// 		perms.HomeDomain = homeDomain
// 	}()

// 	wg.Wait()
// 	return perms, nil
// }
