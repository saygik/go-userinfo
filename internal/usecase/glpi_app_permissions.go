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
	roles, err := u.repo.GetRoles(upn)
	if err == nil {
		perms.Roles = roles
		for _, role := range roles {
			switch role {
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
	// 7. Разделы сайта из user_permissions
	sections, err := u.repo.GetSections(upn)
	if err == nil {
		perms.Sections = appendUnique(perms.Sections, sections...)
	}

	// 8. Авто-добавление AD разделов
	if len(perms.AdminDomains) > 0 || len(perms.TechDomains) > 0 || len(perms.UserDomains) > 0 {
		perms.Sections = appendUnique(perms.Sections, "/ad/users")
	}
	if len(perms.TechDomains) > 0 || len(perms.AdminDomains) > 0 {
		perms.Sections = appendUnique(perms.Sections,
			"/ad/computers", "/ad/lastcomputers", "/ad")
	}
	return perms, nil
}

func appendUnique(slice []string, items ...string) []string {
	m := make(map[string]bool)
	for _, s := range slice {
		m[s] = true
	}
	for _, item := range items {
		if !m[item] {
			slice = append(slice, item)
			m[item] = true
		}
	}
	return slice
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
