package usecase

import (
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/saygik/go-userinfo/internal/entity"
)

// isServiceAccountUser сообщает, входит ли пользователь хотя бы в одну из
// сервисных групп. В memberOf AD-клиент сохраняет только CN групп
// (см. pkg/ad-client firstMembersOfCommaStrings), поэтому из настроенного
// distinguishedName группы извлекается CN. Полное значение DN тоже
// сравнивается — на случай, если memberOf начнёт хранить DN целиком.
func isServiceAccountUser(memberOf []string, serviceGroups []string) bool {
	if len(memberOf) == 0 || len(serviceGroups) == 0 {
		return false
	}

	memberGroups := make(map[string]struct{}, len(memberOf))
	for _, group := range memberOf {
		if normalized := normalizeServiceGroupName(group); normalized != "" {
			memberGroups[normalized] = struct{}{}
		}
	}

	for _, groupDN := range serviceGroups {
		if normalized := normalizeServiceGroupName(groupDN); normalized != "" {
			if _, ok := memberGroups[normalized]; ok {
				return true
			}
		}
		if cn := serviceGroupCN(groupDN); cn != "" {
			if _, ok := memberGroups[normalizeServiceGroupName(cn)]; ok {
				return true
			}
		}
	}

	return false
}

// serviceGroupCN возвращает значение первого RDN, если это CN.
func serviceGroupCN(groupDN string) string {
	dn, err := ldap.ParseDN(groupDN)
	if err != nil || len(dn.RDNs) == 0 || len(dn.RDNs[0].Attributes) == 0 {
		return ""
	}
	attr := dn.RDNs[0].Attributes[0]
	if !strings.EqualFold(attr.Type, "CN") {
		return ""
	}
	return attr.Value
}

func normalizeServiceGroupName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// serviceMailGroupName — группа, членство в которой помечает пользователя как
// владельца служебной почты. В отличие от сервисных учётных записей такие
// пользователи не скрываются из списков.
const serviceMailGroupName = "Служебная почта"

// isServiceMailUser сообщает, входит ли пользователь в группу «Служебная почта».
// memberOf может содержать как CN, так и полный DN — поддерживаются оба вида.
func isServiceMailUser(memberOf []string) bool {
	for _, group := range memberOf {
		if strings.EqualFold(strings.TrimSpace(group), serviceMailGroupName) {
			return true
		}
		if cn := serviceGroupCN(group); cn != "" && strings.EqualFold(strings.TrimSpace(cn), serviceMailGroupName) {
			return true
		}
	}
	return false
}

// userDomainName достаёт имя домена из пользователя, сохранённого в Redis:
// при импорте в поле domain записывается entity.DomainList.
func userDomainName(user map[string]interface{}) string {
	if domain, ok := user["domain"].(map[string]interface{}); ok {
		if name, ok := domain["name"].(string); ok {
			return name
		}
	}
	if name, ok := user["domain"].(string); ok {
		return name
	}
	return ""
}

// canSeeServiceAccounts сообщает, доступны ли пользователю сервисные учётные
// записи домена: системный администратор видит их во всех доменах,
// администратор домена — только в своих доменах.
func (u *UseCase) canSeeServiceAccounts(perms *entity.Permissions, domain string) bool {
	if perms == nil {
		return false
	}
	if perms.IsSysAdmin {
		return true
	}
	return hasDomainAccess(perms.AdminDomains, domain)
}
