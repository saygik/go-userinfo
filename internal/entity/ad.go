package entity

import "time"

type DomainList struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type AvatarForm struct {
	Avatar string `form:"avatar" json:"avatar"`
}

type ADInternetGroups struct {
	WhiteList []string `json:"whitelist"`
	Full      []string `json:"full"`
	Tech      []string `json:"tech"`
}

type ADInternetGroupsDN struct {
	WhiteList string `json:"whitelist"`
	Full      string `json:"full"`
	Tech      string `json:"tech"`
}

// TemporaryGroupChange хранит информацию о временном изменении группы пользователя
type TemporaryGroupChange struct {
	User          string    `json:"user"`          // UPN пользователя
	Domain        string    `json:"domain"`        // Домен пользователя
	UserDN        string    `json:"userDN"`        // Distinguished Name пользователя
	PreviousGroup string    `json:"previousGroup"` // Предыдущая группа (whitelist, full, tech, или "")
	NewGroup      string    `json:"newGroup"`      // Новая группа
	CreatedAt     time.Time `json:"createdAt"`     // Время создания временного изменения
	ExpiresAt     time.Time `json:"expiresAt"`     // Время истечения (когда нужно вернуть обратно)
	ChangedBy     string    `json:"changedBy"`     // Кто сделал изменение
}

// GetRedisKey возвращает ключ Redis для хранения временного изменения
func (t *TemporaryGroupChange) GetRedisKey() string {
	return "temp_group_change:" + t.User
}
