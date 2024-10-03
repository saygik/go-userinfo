package entity

type MattermostUser struct {
	Id          string `json:"id"`
	Name        string `json:"username"`
	AuthService string `json:"auth_service,omitempty"`
	AD          string `json:"ad,omitempty"`
	Email       string `json:"email"`
	Nickname    string `json:"nickname,omitempty"`
	Roles       string `json:"roles,omitempty"`
	IsBot       bool   `json:"is_bot,omitempty"`
	Registred   bool   `json:"registred"`
}

// msgp StringMap
type StringMap map[string]string

// Session contains the user session details.
// This struct's serializer methods are auto-generated. If a new field is added/removed,
// please run make gen-serialized.
type MattermostSession struct {
	Id             string    `json:"id"`
	Token          string    `json:"token"`
	CreateAt       int64     `json:"create_at"`
	ExpiresAt      int64     `json:"expires_at"`
	LastActivityAt int64     `json:"last_activity_at"`
	UserId         string    `json:"user_id"`
	DeviceId       string    `json:"device_id"`
	Roles          string    `json:"roles"`
	IsOAuth        bool      `json:"is_oauth"`
	ExpiredNotify  bool      `json:"expired_notify"`
	Props          StringMap `json:"props"`
	Local          bool      `json:"local" db:"-"`
}

type MattermostUserWithSessions struct {
	Id                 string              `json:"id"`
	CreateAt           int64               `json:"create_at,omitempty"`
	UpdateAt           int64               `json:"update_at,omitempty"`
	DeleteAt           int64               `json:"delete_at"`
	Username           string              `json:"username"`
	Password           string              `json:"password,omitempty"`
	AuthData           *string             `json:"auth_data,omitempty"`
	AuthService        string              `json:"auth_service"`
	Email              string              `json:"email"`
	EmailVerified      bool                `json:"email_verified,omitempty"`
	Nickname           string              `json:"nickname"`
	FirstName          string              `json:"first_name"`
	LastName           string              `json:"last_name"`
	Position           string              `json:"position"`
	Roles              string              `json:"roles"`
	AllowMarketing     bool                `json:"allow_marketing,omitempty"`
	LastPasswordUpdate int64               `json:"last_password_update,omitempty"`
	FailedAttempts     int                 `json:"failed_attempts,omitempty"`
	MfaActive          bool                `json:"mfa_active,omitempty"`
	RemoteId           *string             `json:"remote_id,omitempty"`
	LastActivityAt     int64               `json:"last_activity_at,omitempty"`
	IsBot              bool                `json:"is_bot,omitempty"`
	Sessions           []MattermostSession `json:"sessions,omitempty"`
}

type MattermostHrpPost struct {
	Id      int
	FIO     string
	Dolg    string
	Company string
	Mero    string
}
