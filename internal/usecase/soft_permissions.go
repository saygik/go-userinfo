package usecase

import (
	"fmt"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

// normalizeSoftPermissionDate приводит пустую/техническую дату к пустой строке.
func normalizeSoftPermissionDate(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) >= 4 && (value[0:4] == "1900" || value[0:4] == "0001") {
		return ""
	}
	return value
}

// normalizeSoftUserPermissions гарантирует, что у пользователя есть основное
// полномочие, и что первое полномочие всегда называется "Основные полномочия".
func normalizeSoftUserPermissions(perms []entity.SoftUserPermission, isNew bool) []entity.SoftUserPermission {
	result := make([]entity.SoftUserPermission, 0, len(perms)+1)

	for i, permission := range perms {
		name := strings.TrimSpace(permission.Name)
		if i == 0 {
			name = entity.SoftUserMainPermissionName
		}
		if name == "" {
			name = fmt.Sprintf("Дополнительные полномочия %d", i)
		}

		permission.Name = name
		permission.EndDate = normalizeSoftPermissionDate(permission.EndDate)
		if permission.EndDate == "" {
			permission.Sended = false
		}
		if isNew {
			permission.Id = 0
			permission.IdUserSoft = 0
		}
		result = append(result, permission)
	}

	if len(result) == 0 {
		result = append(result, entity.SoftUserPermission{Name: entity.SoftUserMainPermissionName})
	}

	return result
}

// softUserResponse формирует ответ API с данными AD и списком полномочий.
func (u *UseCase) softUserResponse(user entity.SoftUser) map[string]interface{} {
	adUser := u.GetUserADPropertysShort(user.Name)
	adUser["id"] = user.Id
	adUser["id_soft"] = user.IdSoft
	adUser["login"] = user.Login
	adUser["comment"] = user.Comment
	adUser["fio"] = user.Fio
	adUser["external"] = user.External
	adUser["editor"] = user.Editor
	adUser["last_changes"] = user.LastChanges

	if len(user.Mail) > 0 {
		adUser["mail"] = user.Mail
	} else {
		adUser["admail"] = adUser["mail"]
		adUser["mail"] = ""
	}

	permissions := make([]entity.SoftUserPermission, 0, len(user.Permissions))
	for _, permission := range user.Permissions {
		permission.EndDate = normalizeSoftPermissionDate(permission.EndDate)
		if permission.EndDate == "" {
			permission.Sended = false
		}
		permissions = append(permissions, permission)
	}
	adUser["permissions"] = permissions

	// Совместимость со старым API: срок основного полномочия.
	if len(permissions) > 0 {
		adUser["enddate"] = permissions[0].EndDate
		adUser["sended"] = permissions[0].Sended
	} else {
		adUser["enddate"] = ""
		adUser["sended"] = false
	}

	return adUser
}
