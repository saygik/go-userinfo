package app

import "github.com/saygik/go-userinfo/pkg/logger"

func (a *App) initLogger(env string) {
	l := logger.NewLogger(env, "api.log")

	a.log = l
}
