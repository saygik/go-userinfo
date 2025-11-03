package app

import (
	"context"
	"time"

	"github.com/saygik/go-userinfo/internal/config"
	"github.com/sirupsen/logrus"
)

type App struct {
	cfg config.Config
	c   *Container
	log *logrus.Logger
}

func New(ctx context.Context) (*App, error) {
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

	c := NewAppContainer(ctx, msSQLConnect, glpiConnect, redisConnect, adClients, adConfigs, mattClient, glpiApiClient, hydraClient, oAuth2Client, oAuth2ClientAuthentik, app.cfg.ApiIntegrations.N8nWebhookIvc2Kaspersky, app.log)
	app.c = c
	app.c.useCase.ClearRedisCaсhe()

	//FOR TEST!!!!!!!!!!!!!!!!!!!!
	//	app.c.useCase.GetScheduleTasksNotifications()
	//	duration := 20 * time.Second
	//	time.Sleep(duration)
	//	app.c.useCase.GetHRPTickets()
	//app.c.useCase.GetSoftwareUsersEOL()
	go app.c.useCase.FillRedisCaсheFromAD()
	go app.fillRedis(ctx)
	if app.cfg.App.Env == "prod" {
		go app.getHrpTickets(ctx)
	}
	go app.getCalendarTaskNotifikations(ctx)
	go app.getSoftwareUsersEOL(ctx)

	return app, nil
}

func (a *App) fillRedis(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.c.useCase.FillRedisCaсheFromAD()
		}
	}
}

func (a *App) getHrpTickets(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.c.useCase.GetHRPTickets()
		}
	}
}

func (a *App) getCalendarTaskNotifikations(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.c.useCase.GetScheduleTasksNotifications()
		}
	}
}

// ** Проверка срока окончания действия учетных записей в системах, оповещение пользователя по emnail, создание задачи календаря администраторам системы на отключение
func (a *App) getSoftwareUsersEOL(ctx context.Context) {
	softwarebottime := a.cfg.App.Softwarebottime
	if softwarebottime == 0 {
		softwarebottime = 10
	}
	ticker := time.NewTicker(time.Duration(softwarebottime) * time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// do stuff
			a.c.useCase.GetSoftwareUsersEOL()
		}
	}
}
