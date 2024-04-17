package config

import (
	"time"

	"github.com/ory/viper"
)

// type DBConfig struct {
// 	DBServer   string `json:"CLCK_DB_SERVER"`
// 	DBName     string `json:"CLCK_DB_NAME"`
// 	DBUser     string `json:"CLCK_DB_USER"`
// 	DBPassword string `json:"CLCK_DB_PASS"`
// }

type ADConfig struct {
	Key            string        `json:"key"`
	Name           string        `json:"name"`
	Base           string        `json:"base"`
	Dc             string        `json:"dc"`
	GroupFilter    string        `json:"groupFilter"`
	Filter         string        `json:"filter"`
	ComputerFilter string        `json:"computerFilter"`
	BindDN         string        `json:"bindDN"`
	BindPassword   string        `json:"bind-password"`
	Time           time.Duration `json:"time"`
}

type AppConfig struct {
	Env           string `json:"env"`
	Port          string `json:"port"`
	DefaultDomain string `json:"default-domain"`
}
type DBConfig struct {
	Server   string `json:"server" binding:"required"`
	Dbname   string `json:"dbname" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
	Secret   string
}
type APIConfig struct {
	Server    string `json:"server" binding:"required"`
	Token     string `json:"token" binding:"required"`
	UserToken string `json:"usertoken"`
}
type JWT struct {
	AccessSecret  string `json:"access-secret" binding:"required"`
	RefreshSecret string `json:"refreshsecret" binding:"required"`
	AtExpires     int    `json:"atexpires" binding:"required"`
	RtExpires     int    `json:"rtexpires" binding:"required"`
}
type Repository struct {
	Mssql      DBConfig
	Redis      DBConfig
	Glpi       DBConfig
	Mattermost APIConfig
	GlpiApi    APIConfig
}
type Config struct {
	App        AppConfig `json:"app" binding:"required"`
	Jwt        JWT
	AD         []ADConfig `json:"ad" binding:"required"`
	Repository Repository `json:"repository" binding:"required"`
}

func NewConfig(filePath string) (Config, error) {
	conf := Config{}
	viper.SetConfigType("json")
	//	viper.SetConfigFile(filePath)
	viper.SetConfigName(filePath)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return conf, err
	}
	if err := viper.Unmarshal(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}

// func Load() (cfg Config, err error) {
// 	err = godotenv.Load()
// 	if err != nil {
// 		return cfg, err
// 	}
// 	cfg.Env = getEnv("ENV", "local")
// 	cfg.Port = getEnv("PORT", "8080")
// 	return cfg, nil
// }

// func getEnv(key string, defaultVal string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}

// 	return defaultVal
// }
