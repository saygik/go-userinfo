package usecase

import "github.com/saygik/go-userinfo/internal/entity"

func (u *UseCase) UserInGropScopes(user string, scopes []string, idpScopes []entity.IDPScope) error {
	groupScopes := filterStringArrayByWord(scopes, "group")
	// allowScopes := true
	if len(groupScopes) < 1 {
		return nil
	}
	domain := getDomainFromUserName(user)
	if !u.ad.IsDomainExist(domain) {
		return u.Error("домен пользователя отсутствует в системе")
	}
	groupScopesAD := filterStringArrayByWord(groupScopes, "group-ad")
	for _, groupScope := range groupScopesAD {
		for _, idpScope := range idpScopes {
			if idpScope.Scope == groupScope && idpScope.Domain == domain {
				err := u.UserInDomainGroup(user, idpScope.Group)
				if err == nil {
					return nil
				}
			}
		}
	}

	return u.Error("отсутствуют права пользователя на данную систему")
}
