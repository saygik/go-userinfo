package config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/ory/viper"
)

//	type DBConfig struct {
//		DBServer   string `json:"CLCK_DB_SERVER"`
//		DBName     string `json:"CLCK_DB_NAME"`
//		DBUser     string `json:"CLCK_DB_USER"`
//		DBPassword string `json:"CLCK_DB_PASS"`
//	}
type ADInternetGroups struct {
	WhiteList string `json:"whitelist"`
	Full      string `json:"full"`
	Tech      string `json:"tech"`
}

type ADConfig struct {
	Key            string           `json:"key"`
	Name           string           `json:"name"`
	Base           string           `json:"base"`
	Dc             string           `json:"dc"`
	GroupFilter    string           `json:"groupFilter"`
	Filter         string           `json:"filter"`
	ComputerFilter string           `json:"computerFilter"`
	BindDN         string           `json:"bindDN"`
	BindPassword   string           `json:"bindpassword"`
	AdminGLPIGroup int              `json:"adminGLPIGroup"`
	RmsPort        int              `json:"rmsPort"`
	InternetGroups ADInternetGroups `json:"internetGroups"`
	Time           time.Duration    `json:"time"`
}

type AppConfig struct {
	Env             string `json:"env"`
	Port            string `json:"port"`
	Softwarebottime int    `json:"softwarebottime"`
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

type Repository struct {
	Mssql      DBConfig
	Redis      DBConfig
	Glpi       DBConfig
	Mattermost APIConfig
	GlpiApi    APIConfig
}

type Integrations struct {
	AddCommentFromApi                  string   `json:"add-comment-from-api,omitempty"`
	DisableCalendarTaskNotificationApi string   `json:"disable-calendar-task-notification-api,omitempty"`
	AllowedHosts                       []string `json:"allowed-hosts,omitempty"`
	N8nWebhookIvc2Kaspersky            string   `json:"n8n-webhook-ivc2-kaspersky"`
}

type VaultConfig struct {
	Server     string `json:"server" binding:"required"`
	RoleId     string `json:"roleid" binding:"required"`
	SecretId   string `json:"secretid" binding:"required"`
	SecretPath string `json:"secretpath" binding:"required"`
}

type AuthentikConfig struct {
	Url          string `json:"url" binding:"required"`
	ClientId     string `json:"client-id" binding:"required"`
	ClientSecret string `json:"client-secret" binding:"required"`
	RedirectUrl  string `json:"redirect-url" binding:"required"`
	Scopes       []string
	LogOutUrl    string `json:"log-out-url" binding:"required"`
}
type Config struct {
	App             AppConfig
	Vault           VaultConfig
	AD              []ADConfig
	Repository      Repository
	ApiIntegrations Integrations
	Authentik       AuthentikConfig
}

func NewConfig(filePath string) (Config, error) {
	conf := Config{}
	viper.SetConfigType("json")
	//	viper.SetConfigFile(filePath)
	viper.SetConfigName(filePath)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	if err := viper.Unmarshal(&conf); err != nil {
		return Config{}, err
	}
	config, err := vaultConfig(conf)
	if err != nil {
		return Config{}, err
	}

	return *config, nil
}

func vaultConfig(conf Config) (*Config, error) {
	cfg := &Config{}

	ctx := context.Background()
	cl, err := initVaultClient()
	if err != nil {
		return cfg, err
	}
	resp, err := cl.Auth.AppRoleLogin(
		ctx,
		schema.AppRoleLoginRequest{
			RoleId:   conf.Vault.RoleId,
			SecretId: conf.Vault.SecretId,
		},
	)
	if err != nil {
		return cfg, err
	}
	if err := cl.SetToken(resp.Auth.ClientToken); err != nil {
		return cfg, err
	}
	secret, err := cl.Read(context.Background(), fmt.Sprintf("%sad", conf.Vault.SecretPath))
	if err != nil {
		return cfg, fmt.Errorf("ошибка Vault: %v", err)
	}
	value, ok := secret.Data["data"]
	if !ok {
		return cfg, fmt.Errorf("ошибка Vault: %v", err)
	}
	arrOfinterface, ok := value.(map[string]interface{})["domains"].([]interface{})
	if !ok {
		return cfg, fmt.Errorf("ошибка Vault: %v", err)
	}
	data, err := json.Marshal(arrOfinterface)
	if err != nil {
		return cfg, fmt.Errorf("ошибка Vault: %v", err)
	}
	adOne := []ADConfig{}
	err = json.Unmarshal(data, &adOne)
	if err != nil {
		return cfg, fmt.Errorf("ошибка Vault: %v", err)
	}
	cfg.AD = adOne

	secret, err = cl.Read(context.Background(), fmt.Sprintf("%srepository", conf.Vault.SecretPath))
	if err != nil {
		return cfg, fmt.Errorf("ошибка Vault: %v", err)
	}
	value, ok = secret.Data["data"]
	if !ok {
		return cfg, fmt.Errorf("ошибка Vault: %v", err)
	}

	data, err = json.Marshal(value)
	if err != nil {
		return cfg, fmt.Errorf("ошибка Vault: %v", err)
	}
	repo := Repository{}
	err = json.Unmarshal(data, &repo)
	if err != nil {
		return cfg, fmt.Errorf("ошибка Vault: %v", err)
	}
	cfg.Repository = repo

	//* authentik
	secret, err = cl.Read(context.Background(), fmt.Sprintf("%sauthentik", conf.Vault.SecretPath))
	if err != nil {
		return cfg, fmt.Errorf("ошибка Authentik client, секрет по пути %s недоступен: %v", fmt.Sprintf("%sjwt", conf.Vault.SecretPath), err)
	}
	value, ok = secret.Data["data"]
	if !ok {
		return cfg, fmt.Errorf("ошибка Authentik: %v", err)
	}

	data, err = json.Marshal(value)
	if err != nil {
		return cfg, fmt.Errorf("ошибка Authentik: %v", err)
	}
	authentik := AuthentikConfig{}
	err = json.Unmarshal(data, &authentik)
	if err != nil {
		return cfg, fmt.Errorf("ошибка Authentik: %v", err)
	}
	cfg.Authentik = authentik
	//* integrations
	secret, err = cl.Read(context.Background(), fmt.Sprintf("%sintegrations", conf.Vault.SecretPath))
	if err != nil {
		return cfg, fmt.Errorf("ошибка integrations , секрет по пути %s недоступен: %v", fmt.Sprintf("%sjwt", conf.Vault.SecretPath), err)
	}
	value, ok = secret.Data["data"]
	if !ok {
		return cfg, fmt.Errorf("ошибка integrations: %v", err)
	}

	data, err = json.Marshal(value)
	if err != nil {
		return cfg, fmt.Errorf("ошибка Hydra: %v", err)
	}
	integrations := Integrations{}
	err = json.Unmarshal(data, &integrations)
	if err != nil {
		return cfg, fmt.Errorf("ошибка integrations: %v", err)
	}
	cfg.ApiIntegrations = integrations

	_ = resp
	cfg.App = conf.App

	// cfg.AD = conf.AD
	// cfg.Repository = conf.Repository
	// cfg.Jwt = conf.Jwt
	return cfg, nil
}

func initVaultClient() (*vault.Client, error) {
	// prepare a client with the given base address
	client, err := vault.New(
		vault.WithAddress("https://vault.brnv.rw:8200/"),
		vault.WithRequestTimeout(10*time.Second),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
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
