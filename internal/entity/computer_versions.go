package entity

// ComputerVersionCount — агрегированная статистика по версиям ОС компьютеров домена (из AD).
type ComputerVersionCount struct {
	OperatingSystem             string `json:"operatingSystem"`
	OperatingSystemFamily       string `json:"operatingSystemFamily,omitempty"`
	OperatingSystemVersion      string `json:"operatingSystemVersion"`
	OperatingSystemVersionHuman string `json:"operatingSystemVersionHuman,omitempty"`
	Count                       int    `json:"count"`
}

// ComputerFamilyCount — агрегированная статистика по семействам ОС (Windows 11, Windows 10, Windows Server 2012 R2 и т.п.).
type ComputerFamilyCount struct {
	OperatingSystemFamily string `json:"operatingSystemFamily"`
	Count                 int    `json:"count"`
}

