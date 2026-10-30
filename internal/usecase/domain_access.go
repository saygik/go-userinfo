package usecase

import "strings"

// normalizeDomainName приводит имя домена к единому виду: без пробелов и в
// нижнем регистре. Ключи карт прав (AdminDomains, TechDomains, UserDomains)
// хранятся нормализованными, поэтому сравнивать их нужно так же.
func normalizeDomainName(domain string) string {
	return strings.ToLower(strings.TrimSpace(domain))
}

// hasDomainAccess проверяет наличие доступа к домену в карте прав.
// Ключи карт из LoadUserPermissions нормализованы, но helper терпим и к
// «сырым» ключам — сравнивает имена без учёта регистра и пробелов.
func hasDomainAccess(domains map[string]bool, domain string) bool {
	if len(domains) == 0 {
		return false
	}
	if domains[domain] {
		return true
	}
	normalized := normalizeDomainName(domain)
	if domains[normalized] {
		return true
	}
	for key := range domains {
		if normalizeDomainName(key) == normalized {
			return true
		}
	}
	return false
}
