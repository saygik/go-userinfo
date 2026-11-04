package entity

import "time"

// PermissionAudit — запись аудита изменения прав пользователя.
type PermissionAudit struct {
	Id        int64     `db:"id" json:"id"`
	ChangedAt time.Time `db:"changed_at" json:"changed_at"`
	Actor     string    `db:"actor_upn" json:"actor"`
	Target    string    `db:"target_upn" json:"target"`
	Action    string    `db:"action" json:"action"`
	Object    string    `db:"object_name" json:"object,omitempty"`
	Details   string    `db:"details" json:"details,omitempty"`
}
