package usecase

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/saygik/go-userinfo/internal/entity"
	"github.com/saygik/go-userinfo/internal/state"
)

type adminUserAggregate struct {
	upn      string
	roles    map[int]entity.IdNameDescription
	access   map[string]entity.DomainAccess
	sections map[string]entity.IdNameDescription
}

type adminADUser struct {
	fio      string
	domain   string
	disabled bool
	service  bool
}

// GetAdminOverview собирает системную информацию и список пользователей,
// чьи права отличаются от стандартных (только роль DomainUser и доступ
// пользователя к своему домашнему домену).
func (u *UseCase) GetAdminOverview() (entity.AdminOverview, error) {
	warnings := make([]string, 0)

	roleRows, err := u.repo.GetAllUserRoles()
	if err != nil {
		warnings = append(warnings, "не удалось получить роли пользователей: "+err.Error())
	}
	accessRows, err := u.repo.GetAllUserDomainAccess()
	if err != nil {
		warnings = append(warnings, "не удалось получить доступы к доменам: "+err.Error())
	}
	sectionRows, err := u.repo.GetAllUserSections()
	if err != nil {
		warnings = append(warnings, "не удалось получить разделы пользователей: "+err.Error())
	}

	aggregates := make(map[string]*adminUserAggregate)
	roleUsers := make(map[string]bool)
	accessUsers := make(map[string]bool)
	sectionUsers := make(map[string]bool)

	getAggregate := func(upn string) *adminUserAggregate {
		key := strings.ToLower(strings.TrimSpace(upn))
		agg, ok := aggregates[key]
		if !ok {
			agg = &adminUserAggregate{
				upn:      upn,
				roles:    make(map[int]entity.IdNameDescription),
				access:   make(map[string]entity.DomainAccess),
				sections: make(map[string]entity.IdNameDescription),
			}
			aggregates[key] = agg
		}
		return agg
	}

	for _, row := range roleRows {
		if strings.TrimSpace(row.UserPrincipalName) == "" {
			continue
		}
		agg := getAggregate(row.UserPrincipalName)
		agg.roles[row.RoleID] = entity.IdNameDescription{
			Id:          row.RoleID,
			Name:        row.Name,
			Description: row.Description,
		}
		roleUsers[strings.ToLower(strings.TrimSpace(row.UserPrincipalName))] = true
	}

	for _, row := range accessRows {
		if strings.TrimSpace(row.UserPrincipalName) == "" {
			continue
		}
		agg := getAggregate(row.UserPrincipalName)
		key := normalizeDomainName(row.Domain) + "|" + strings.ToLower(strings.TrimSpace(row.AccessLevel))
		agg.access[key] = entity.DomainAccess{Domain: row.Domain, AccessLevel: row.AccessLevel}
		accessUsers[strings.ToLower(strings.TrimSpace(row.UserPrincipalName))] = true
	}

	for _, row := range sectionRows {
		if strings.TrimSpace(row.UserPrincipalName) == "" {
			continue
		}
		agg := getAggregate(row.UserPrincipalName)
		key := strconv.Itoa(row.Id) + "|" + row.Name
		agg.sections[key] = entity.IdNameDescription{
			Id:          row.Id,
			Name:        row.Name,
			Description: row.Description,
		}
		sectionUsers[strings.ToLower(strings.TrimSpace(row.UserPrincipalName))] = true
	}

	adUsers := u.loadAdminADUsers()

	nonStandard := make([]entity.AdminUserPermission, 0, len(aggregates))
	for _, agg := range aggregates {
		user, ok := buildAdminUserPermission(agg, adUsers)
		if ok {
			nonStandard = append(nonStandard, user)
		}
	}

	sort.Slice(nonStandard, func(i, j int) bool {
		return strings.ToLower(nonStandard[i].User) < strings.ToLower(nonStandard[j].User)
	})

	domains := u.ad.DomainList()
	domainNames := make([]string, 0, len(domains))
	for _, domain := range domains {
		domainNames = append(domainNames, domain.Name)
	}
	sort.Strings(domainNames)

	lastFill := ""
	if filledAt := state.LastADFill(); !filledAt.IsZero() {
		lastFill = filledAt.Format(time.RFC3339)
	}

	adUsersCount, adComputersCount, err := u.GetAdCounts()
	if err != nil {
		warnings = append(warnings, "не удалось получить счётчики AD: "+err.Error())
	}
	appRoles, err := u.repo.GetAppRoles()
	if err != nil {
		warnings = append(warnings, "не удалось получить роли приложения: "+err.Error())
	}
	appSections, err := u.repo.GetAppSections()
	if err != nil {
		warnings = append(warnings, "не удалось получить разделы приложения: "+err.Error())
	}
	softwaresCount, err := u.repo.CountUserSoftwares()
	if err != nil {
		warnings = append(warnings, "не удалось получить количество записей пользователей в системах: "+err.Error())
	}
	if appRoles == nil {
		appRoles = []entity.IdNameDescription{}
	}
	if appSections == nil {
		appSections = []entity.IdNameDescription{}
	}

	overview := entity.AdminOverview{
		System: entity.AdminSystemInfo{
			Warnings: warnings,
			Database: pingAdminService(u.repo.Ping),
			Redis:    pingAdminService(u.redis.Ping),
			AD: entity.AdminADInfo{
				Domains:          domainNames,
				DomainCount:      len(domainNames),
				CacheInitialized: state.IsInitialized(),
				CacheFilling:     state.IsFillingRedis(),
				LastFillAt:       lastFill,
			},
			Counts: entity.AdminCounts{
				ADUsers:               adUsersCount,
				ADComputers:           adComputersCount,
				AppRoles:              len(appRoles),
				AppSections:           len(appSections),
				UsersWithRoles:        len(roleUsers),
				UsersWithDomainAccess: len(accessUsers),
				UsersWithSections:     len(sectionUsers),
				UserSoftwares:         softwaresCount,
			},
			Roles:    appRoles,
			Sections: appSections,
		},
		Users: nonStandard,
	}

	return overview, nil
}

func pingAdminService(ping func() error) entity.AdminServiceInfo {
	start := time.Now()
	err := ping()
	info := entity.AdminServiceInfo{
		Available: err == nil,
		LatencyMS: time.Since(start).Milliseconds(),
	}
	if err != nil {
		info.Error = err.Error()
	}
	return info
}

func (u *UseCase) loadAdminADUsers() map[string]adminADUser {
	result := make(map[string]adminADUser)

	redisUsers, err := u.redis.GetKeyFieldAll("allusers")
	if err != nil {
		return result
	}

	for key, raw := range redisUsers {
		var user map[string]interface{}
		if json.Unmarshal([]byte(raw), &user) != nil {
			continue
		}

		upn, _ := user["userPrincipalName"].(string)
		if upn == "" {
			upn = key
		}
		if upn == "" {
			continue
		}

		info := adminADUser{}
		info.fio, _ = user["cn"].(string)
		info.disabled, _ = user["disabled"].(bool)
		info.service, _ = user["service"].(bool)

		switch domain := user["domain"].(type) {
		case map[string]interface{}:
			info.domain, _ = domain["name"].(string)
		case string:
			info.domain = domain
		}

		result[strings.ToLower(strings.TrimSpace(upn))] = info
	}

	return result
}

func buildAdminUserPermission(agg *adminUserAggregate, adUsers map[string]adminADUser) (entity.AdminUserPermission, bool) {
	roles := make([]entity.IdNameDescription, 0, len(agg.roles))
	for _, role := range agg.roles {
		roles = append(roles, role)
	}
	sort.Slice(roles, func(i, j int) bool { return roles[i].Id < roles[j].Id })

	access := make([]entity.DomainAccess, 0, len(agg.access))
	for _, item := range agg.access {
		access = append(access, item)
	}
	sort.Slice(access, func(i, j int) bool {
		if normalizeDomainName(access[i].Domain) == normalizeDomainName(access[j].Domain) {
			return strings.ToLower(access[i].AccessLevel) < strings.ToLower(access[j].AccessLevel)
		}
		return normalizeDomainName(access[i].Domain) < normalizeDomainName(access[j].Domain)
	})

	sections := make([]entity.IdNameDescription, 0, len(agg.sections))
	for _, section := range agg.sections {
		sections = append(sections, section)
	}
	sort.Slice(sections, func(i, j int) bool { return sections[i].Name < sections[j].Name })

	user := entity.AdminUserPermission{
		User:         agg.upn,
		Roles:        roles,
		DomainAccess: access,
		Sections:     sections,
		Reasons:      []string{},
	}

	if adUser, found := adUsers[strings.ToLower(strings.TrimSpace(agg.upn))]; found {
		user.FoundInAD = true
		user.Fio = adUser.fio
		user.Domain = adUser.domain
		user.Disabled = adUser.disabled
		user.Service = adUser.service
	} else {
		user.Domain = normalizeDomainName(getDomainFromUserName(agg.upn))
	}

	hasAllDomains := false
	extraRoles := make([]string, 0, len(roles))
	for _, role := range roles {
		switch role.Name {
		case "DomainUser":
			hasAllDomains = true
		case "SysAdmin":
			user.IsSysAdmin = true
			user.IsAdmin = true
			user.IsTech = true
			extraRoles = append(extraRoles, role.Name)
		case "Admin":
			user.IsAdmin = true
			user.IsTech = true
			extraRoles = append(extraRoles, role.Name)
		case "Tech":
			user.IsTech = true
			extraRoles = append(extraRoles, role.Name)
		default:
			extraRoles = append(extraRoles, role.Name)
		}
	}
	if len(extraRoles) > 0 {
		user.Reasons = append(user.Reasons, "роли: "+strings.Join(extraRoles, ", "))
	}

	homeDomain := normalizeDomainName(getDomainFromUserName(agg.upn))
	var adminDomains, techDomains, otherUserDomains []string
	for _, item := range access {
		domain := normalizeDomainName(item.Domain)
		switch strings.ToLower(strings.TrimSpace(item.AccessLevel)) {
		case "admin":
			adminDomains = append(adminDomains, domain)
		case "tech":
			techDomains = append(techDomains, domain)
		case "user":
			// Роль DomainUser уже даёт доступ пользователя ко всем доменам,
			// поэтому явные записи user для неё не считаются расширением прав.
			if !hasAllDomains && domain != homeDomain {
				otherUserDomains = append(otherUserDomains, domain)
			}
		}
	}
	if len(adminDomains) > 0 {
		user.Reasons = append(user.Reasons, "администратор домена: "+strings.Join(adminDomains, ", "))
	}
	if len(techDomains) > 0 {
		user.Reasons = append(user.Reasons, "технический доступ: "+strings.Join(techDomains, ", "))
	}
	if len(otherUserDomains) > 0 {
		user.Reasons = append(user.Reasons, "доступ пользователя к чужим доменам: "+strings.Join(otherUserDomains, ", "))
	}

	if len(sections) > 0 {
		names := make([]string, 0, len(sections))
		for _, section := range sections {
			name := section.Name
			if name == "" {
				name = section.Description
			}
			names = append(names, name)
		}
		sort.Strings(names)
		user.Reasons = append(user.Reasons, "дополнительные разделы: "+strings.Join(names, ", "))
	}

	return user, len(user.Reasons) > 0
}
