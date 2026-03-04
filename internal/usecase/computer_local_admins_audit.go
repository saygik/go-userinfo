package usecase

import (
	"encoding/json"
	"strings"

	"github.com/saygik/go-userinfo/internal/entity"
)

func (u *UseCase) ComputerLocalAdminsAudit(computer string, localAdmins []string, isDomain bool) error {
	return u.repo.ComputerLocalAdminsAudit(computer, localAdmins, isDomain)
}

func (u *UseCase) ComputerLocalAdminsGet(isDomain bool) (results []entity.LocalAdmins, err error) {
	return u.repo.ComputerLocalAdminsGet(isDomain)
}

// UpdateComputerLocalAdmins сохраняет запись в БД и обновляет кеш "adc" для указанного компьютера.
func (u *UseCase) UpdateComputerLocalAdmins(perms entity.Permissions, computer, domain, administrators string) error {
	if !HasAdminAccess(perms, domain) {
		return u.Error("у вас нет прав для выполнения операции")
	}
	// 1. Записываем в БД
	if err := u.repo.UpdateComputerLocalAdmins(computer, "1", administrators); err != nil {
		return err
	}

	// 2. Обновляем кеш Redis "adc" для домена
	adcStr, err := u.redis.GetKeyFieldValue("adc", domain)
	if err != nil || adcStr == "" {
		// Если в кеше ничего нет — тихо выходим, БД уже обновлена
		return nil
	}

	var comps []map[string]any
	if err := json.Unmarshal([]byte(adcStr), &comps); err != nil {
		return err
	}

	updated := false
	for i := range comps {
		name, _ := comps[i]["name"].(string)
		if strings.EqualFold(name, computer) {
			comps[i]["administrators_domain"] = administrators
			updated = true
			break
		}
	}
	if !updated {
		// Компьютер не найден в кеше домена
		return nil
	}

	newJSON, err := json.Marshal(comps)
	if err != nil {
		return err
	}

	if err := u.redis.AddKeyFieldValue("adc", domain, newJSON); err != nil {
		return err
	}

	return nil
}

func HasAdminAccess(perms entity.Permissions, domain string) bool {
	// 1. 🔥 IsSysAdmin — имеет доступ ко всем доменам
	if perms.IsSysAdmin {
		return true
	}

	// 2. 🔥 domain в AdminDomains
	if perms.AdminDomains != nil {
		return perms.AdminDomains[domain]
	}

	return false
}
