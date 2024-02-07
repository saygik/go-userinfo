package app

import (
	"github.com/saygik/go-userinfo/internal/controller/http"
	"github.com/saygik/go-userinfo/internal/controller/http/api"
)

func (a *App) StartHTTPServer() error {

	s := http.NewServer(a.cfg.App.Env, a.log)

	api.NewHandler(s.Rtr, a.C.GetUseCase(), a.log)

	err := s.Start(a.cfg.App.Port)
	if err != nil {
		return err
	}

	return nil
}
