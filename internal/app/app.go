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
	//app.log.Info("------------------Starting programm-------------")
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
	adClients := NewADClients(cfg.AD)

	mattClient := app.newMattermostConnection(cfg.Repository.Mattermost.Server, cfg.Repository.Mattermost.Token)
	c := NewAppContainer(msSQLConnect, glpiConnect, redisConnect, adClients, mattClient)
	app.c = c
	app.c.GetUseCase().ClearRedisCaсhe()
	app.c.GetUseCase().FillRedisCaсheFromAD()
	ticker := time.NewTicker(30 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
				app.c.GetUseCase().FillRedisCaсheFromAD()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return app, nil
}
