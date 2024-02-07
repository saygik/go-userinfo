package app

import (
	"github.com/saygik/go-userinfo/internal/config"
	"github.com/sirupsen/logrus"
)

type App struct {
	cfg config.Config
	C   *Container
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
	c := NewAppContainer(msSQLConnect, glpiConnect, redisConnect, adClients)
	app.C = c
	app.C.GetUseCase().FillRedisCaсheFromAD()
	return app, nil
}
