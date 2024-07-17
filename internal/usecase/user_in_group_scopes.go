package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) UserInGropScopes(user string, scopes []string, idpScopes []entity.IDPScope, client *entity.OAuth2Client) (roles []string, useRoles bool, err error) {
	access := false
	useRoles = IsStringInArray("role", scopes)
	if useRoles {
		roles = addStringToArrayIfNotExist("Viewer", roles)
	}
	groupScopes := filterStringArrayByWord(scopes, "group")
	if len(groupScopes) < 1 {
		return roles, useRoles, nil
	}
	domain := getDomainFromUserName(user)
	if !u.ad.IsDomainExist(domain) {
		return roles, false, u.Error("домен пользователя отсутствует в системе")
	}
	groupScopesAD := filterStringArrayByWord(groupScopes, "group-ad")
	for _, groupScope := range groupScopesAD {
		for _, idpScope := range idpScopes {
			if idpScope.Scope == groupScope && idpScope.Domain == domain {
				err := u.UserInDomainGroup(user, idpScope.Group)
				if err == nil {
					if useRoles {
						roles = addStringToArrayIfNotExist(idpScope.Role, roles)
					}
					access = true
				}
			}
		}
	}

	if access {
		return roles, useRoles, nil
	}
	errMsg := "отсутствуют права пользователя на данную систему"
	if client == nil {
		return roles, false, u.Error(errMsg)
	}
	if len(*client.ClientName) > 0 {
		errMsg = "отсутствуют права пользователя на систему " + *client.ClientName
	} else if len(*client.ClientId) > 0 {
		errMsg = "отсутствуют права пользователя на систему " + *client.ClientId
	}
	return roles, false, u.Error(errMsg)
}
