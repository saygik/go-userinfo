package entity

type MattermostUser struct {
	Id          string `json:"id"`
	Name        string `json:"username"`
	AuthService string `json:"auth_service,omitempty"`
	AD          string `json:"ad,omitempty"`
	Email       string `json:"email"`
	Nickname    string `json:"nickname,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
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
	Status             string              `json:"status,omitempty"`
}

type MattermostIntegration struct {
	URL     string `json:"url"`
	Context struct {
		Action string `json:"action"`
	} `json:"context"`
}

type MattermostAction struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Integration MattermostIntegration `json:"integration"`
}

type MattermostAttachment struct {
	Text    string             `json:"text"`
	Actions []MattermostAction `json:"actions"`
}

type User struct {
	Id                     string    `json:"id"`
	CreateAt               int64     `json:"create_at,omitempty"`
	UpdateAt               int64     `json:"update_at,omitempty"`
	DeleteAt               int64     `json:"delete_at"`
	Username               string    `json:"username"`
	Password               string    `json:"password,omitempty"`
	AuthData               *string   `json:"auth_data,omitempty"`
	AuthService            string    `json:"auth_service"`
	Email                  string    `json:"email"`
	EmailVerified          bool      `json:"email_verified,omitempty"`
	Nickname               string    `json:"nickname"`
	FirstName              string    `json:"first_name"`
	LastName               string    `json:"last_name"`
	Position               string    `json:"position"`
	Roles                  string    `json:"roles"`
	AllowMarketing         bool      `json:"allow_marketing,omitempty"`
	Props                  StringMap `json:"props,omitempty"`
	NotifyProps            StringMap `json:"notify_props,omitempty"`
	LastPasswordUpdate     int64     `json:"last_password_update,omitempty"`
	LastPictureUpdate      int64     `json:"last_picture_update,omitempty"`
	FailedAttempts         int       `json:"failed_attempts,omitempty"`
	Locale                 string    `json:"locale"`
	Timezone               StringMap `json:"timezone"`
	MfaActive              bool      `json:"mfa_active,omitempty"`
	MfaSecret              string    `json:"mfa_secret,omitempty"`
	RemoteId               *string   `json:"remote_id,omitempty"`
	LastActivityAt         int64     `json:"last_activity_at,omitempty"`
	IsBot                  bool      `json:"is_bot,omitempty"`
	BotDescription         string    `json:"bot_description,omitempty"`
	BotLastIconUpdate      int64     `json:"bot_last_icon_update,omitempty"`
	TermsOfServiceId       string    `json:"terms_of_service_id,omitempty"`
	TermsOfServiceCreateAt int64     `json:"terms_of_service_create_at,omitempty"`
	DisableWelcomeEmail    bool      `json:"disable_welcome_email"`
	LastLogin              int64     `json:"last_login,omitempty"`
}

type MattermostInteractiveMessageRequestForm struct {
	UserId    string `validate:"required" json:"user_id"`
	PostId    string `validate:"required" json:"post_id"`
	ChannelId string `validate:"required" json:"channel_id"`
	TeamId    string `validate:"required" json:"team_id"`
	Context   struct {
		Comment string `json:"comment,omitempty"`
		Action  string `json:"action,omitempty"`
		Id      int    `json:"id,omitempty"`
		Soft    string `json:"soft,omitempty"`
	} `json:"context"`
}
