package app

import (
	"time"

	"github.com/saygik/go-userinfo/internal/config"
	"github.com/sirupsen/logrus"
)

type App struct {
	cfg config.Config
	c   *Container
	log *logrus.Logger
}

func New() (*App, error) {
	app := &App{}
	cfg, err := config.NewConfig("config.json")
	if err != nil {
		return nil, err
	}
	app.cfg = cfg
	app.initLogger(cfg.App.Env)
	app.log.Info("------------------Starting programm-------------")
	msSQLConnect, err := app.newMsSQLConnect(cfg.Repository.Mssql)
	if err != nil {
		return nil, err
	}
	glpiConnect, err := app.newGLPISQLConnect(cfg.Repository.Glpi)
	if err != nil {
		return nil, err
	}
	redisConnect, err := app.newRedisConnect(cfg.Repository.Redis.Server, cfg.Repository.Redis.Password, 0)
	if err != nil {
		return nil, err
	}
	hydraClient, err := app.newHydraClient(cfg.Hydra.Url, cfg.Hydra.IDPScopes)
	if err != nil {
		return nil, err
	}
	oAuth2Client, err := app.newOAuth2Client(cfg.Hydra.Url, cfg.Hydra.ClientId, cfg.Hydra.ClientSecret, cfg.Hydra.RedirectUrl, cfg.Hydra.Scopes, cfg.Hydra.LogOutUrl)
	if err != nil {
		return nil, err
	}
	oAuth2ClientAuthentik, err := app.newOAuth2Client(cfg.Authentik.Url, cfg.Authentik.ClientId, cfg.Authentik.ClientSecret, cfg.Authentik.RedirectUrl, cfg.Authentik.Scopes, cfg.Authentik.LogOutUrl)
	if err != nil {
		return nil, err
	}

	adClients, adConfigs := NewADClients(cfg.AD)

	mattClient := app.newMattermostConnection(cfg.Repository.Mattermost.Server, cfg.Repository.Mattermost.Token, cfg.ApiIntegrations.AddCommentFromApi, cfg.ApiIntegrations.DisableCalendarTaskNotificationApi, cfg.ApiIntegrations.AllowedHosts)
	glpiApiClient := app.newGLPIApiConnection(cfg.Repository.GlpiApi.Server, cfg.Repository.GlpiApi.Token, cfg.Repository.GlpiApi.UserToken)

	c := NewAppContainer(msSQLConnect, glpiConnect, redisConnect, adClients, adConfigs, mattClient, glpiApiClient, hydraClient, oAuth2Client, oAuth2ClientAuthentik, app.cfg.ApiIntegrations.N8nWebhookIvc2Kaspersky, app.log)
	app.c = c
	app.c.GetUseCase().ClearRedisCaсhe()
	go app.c.GetUseCase().FillRedisCaсheFromAD()
	ticker := time.NewTicker(30 * time.Minute)

	quit := make(chan struct{})
	fillingRedis := false
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				fillingRedis = true
				app.c.GetUseCase().FillRedisCaсheFromAD()
				fillingRedis = false
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	ticker2 := time.NewTicker(1 * time.Minute)
	quit2 := make(chan struct{})

	//FOR TEST!!!!!!!!!!!!!!!!!!!!
	//	app.c.GetUseCase().GetScheduleTasksNotifications()
	//app.c.GetUseCase().GetHRPTickets()
	//app.c.GetUseCase().GetSoftwareUsersEOL()
	if app.cfg.App.Env == "prod" {
		go getHrpTickets(app, ticker2, quit2, &fillingRedis)
	}

	ticker3 := time.NewTicker(1 * time.Minute)
	quit3 := make(chan struct{})
	go getCalendarTaskNotifikations(app, ticker3, quit3)

	//	ticker4 := time.NewTicker(24 * time.Hour)
	softwarebottime := app.cfg.App.Softwarebottime
	if softwarebottime == 0 {
		softwarebottime = 10
	}
	ticker4 := time.NewTicker(time.Duration(softwarebottime) * time.Minute)
	quit4 := make(chan struct{})
	go getSoftwareUsersEOL(app, ticker4, quit4)

	return app, nil
}

func getHrpTickets(app *App, ticker2 *time.Ticker, quit2 chan struct{}, fillingRedis *bool) {
	for {
		select {
		case <-ticker2.C:
			// do stuff
			if !*fillingRedis {
				app.c.GetUseCase().GetHRPTickets()
			}
		case <-quit2:
			ticker2.Stop()
			return
		}
	}
}

func getCalendarTaskNotifikations(app *App, ticker3 *time.Ticker, quit3 chan struct{}) {
	for {
		select {
		case <-ticker3.C:
			// do stuff
			app.c.GetUseCase().GetScheduleTasksNotifications()
		case <-quit3:
			ticker3.Stop()
			return
		}
	}
}

// ** Проверка срока окончания действия учетных записей в системах, оповещение пользователя по emnail, создание задачи календаря администраторам системы на отключение
func getSoftwareUsersEOL(app *App, ticker4 *time.Ticker, quit4 chan struct{}) {
	for {
		select {
		case <-ticker4.C:
			// do stuff
			app.c.GetUseCase().GetSoftwareUsersEOL()
		case <-quit4:
			ticker4.Stop()
			return
		}
	}
}
