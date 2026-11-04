package entity

// AdminOverview — сводная информация для панели администратора.
type AdminOverview struct {
	System AdminSystemInfo       `json:"system"`
	Users  []AdminUserPermission `json:"users"`
}

// AdminServiceInfo описывает состояние внешнего сервиса (БД, Redis).
type AdminServiceInfo struct {
	Available bool   `json:"available"`
	LatencyMS int64  `json:"latency_ms"`
	Error     string `json:"error,omitempty"`
}

// AdminADInfo описывает состояние источников Active Directory и кеша.
type AdminADInfo struct {
	Domains          []string `json:"domains"`
	DomainCount      int      `json:"domain_count"`
	CacheInitialized bool     `json:"cache_initialized"`
	CacheFilling     bool     `json:"cache_filling"`
	LastFillAt       string   `json:"last_fill_at,omitempty"`
}

// AdminCounts — основные счётчики системы.
type AdminCounts struct {
	ADUsers               int `json:"ad_users"`
	ADComputers           int `json:"ad_computers"`
	AppRoles              int `json:"app_roles"`
	AppSections           int `json:"app_sections"`
	UsersWithRoles        int `json:"users_with_roles"`
	UsersWithDomainAccess int `json:"users_with_domain_access"`
	UsersWithSections     int `json:"users_with_sections"`
	UserSoftwares         int `json:"user_softwares"`
}

// AdminSystemInfo — системная информация панели администратора.
type AdminSystemInfo struct {
	Version       string           `json:"version"`
	Commit        string           `json:"commit"`
	Built         string           `json:"built"`
	GoVersion     string           `json:"go_version"`
	Platform      string           `json:"platform"`
	NumCPU        int              `json:"num_cpu"`
	GOMAXPROCS    int              `json:"gomaxprocs"`
	Goroutines    int              `json:"goroutines"`
	MemoryAllocMB float64          `json:"memory_alloc_mb"`
	MemorySysMB   float64          `json:"memory_sys_mb"`
	StartedAt     string           `json:"started_at"`
	ServerTime    string           `json:"server_time"`
	UptimeSeconds int64            `json:"uptime_seconds"`
	Database      AdminServiceInfo `json:"database"`
	Redis         AdminServiceInfo `json:"redis"`
	AD            AdminADInfo      `json:"ad"`
	Counts        AdminCounts      `json:"counts"`
	Warnings      []string         `json:"warnings,omitempty"`
}

// AdminUserPermission — пользователь с правами, отличными от стандартных.
type AdminUserPermission struct {
	User         string              `json:"user"`
	Fio          string              `json:"fio,omitempty"`
	Domain       string              `json:"domain,omitempty"`
	FoundInAD    bool                `json:"found_in_ad"`
	Disabled     bool                `json:"disabled"`
	Service      bool                `json:"service"`
	IsSysAdmin   bool                `json:"is_sysadmin"`
	IsAdmin      bool                `json:"is_admin"`
	IsTech       bool                `json:"is_tech"`
	Roles        []IdNameDescription `json:"roles"`
	DomainAccess []DomainAccess      `json:"domain_access"`
	Sections     []IdNameDescription `json:"sections"`
	Reasons      []string            `json:"reasons"`
}

// AdminUserRoleRow — плоская строка роли пользователя для агрегации.
type AdminUserRoleRow struct {
	UserPrincipalName string `db:"user_principal_name"`
	RoleID            int    `db:"role_id"`
	Name              string `db:"name"`
	Description       string `db:"description"`
}

// AdminUserDomainAccessRow — плоская строка доступа пользователя к домену.
type AdminUserDomainAccessRow struct {
	UserPrincipalName string `db:"user_principal_name"`
	Domain            string `db:"domain"`
	AccessLevel       string `db:"access_level"`
}

// AdminUserSectionRow — плоская строка дополнительного раздела пользователя.
type AdminUserSectionRow struct {
	UserPrincipalName string `db:"user_principal_name"`
	Id                int    `db:"id"`
	Name              string `db:"name"`
	Description       string `db:"description"`
}
