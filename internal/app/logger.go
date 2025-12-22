package app

import "github.com/saygik/go-userinfo/pkg/logger"

func (a *App) initLogger(env string) {
	a.log = logger.NewLogger(env, "app.log")
	a.accessLog = logger.NewLogger(env, "access.log")
}
