package entity

import "time"

// UserSoftwareAudit — запись аудита изменения связи пользователя с системой.
type UserSoftwareAudit struct {
	Id             int64     `db:"id" json:"id"`
	ChangedAt      time.Time `db:"changed_at" json:"changed_at"`
	Actor          string    `db:"actor_upn" json:"actor"`
	Action         string    `db:"action" json:"action"`
	UserSoftwareID int64     `db:"user_software_id" json:"user_software_id,omitempty"`
	IdSoft         int       `db:"id_soft" json:"id_soft,omitempty"`
	SoftName       string    `db:"soft_name" json:"soft_name,omitempty"`
	UserName       string    `db:"user_name" json:"user_name,omitempty"`
	UserLogin      string    `db:"user_login" json:"user_login,omitempty"`
	Fio            string    `db:"fio" json:"fio,omitempty"`
	Details        string    `db:"details" json:"details,omitempty"`
}
