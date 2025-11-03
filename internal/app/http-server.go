package app

import (
	_ "net/http/pprof"

	"github.com/saygik/go-userinfo/internal/controller/http"
	"github.com/saygik/go-userinfo/internal/controller/http/api"
	"github.com/saygik/go-userinfo/internal/controller/metrics"
)

func (a *App) StartHTTPServer() error {

	s := http.NewServer(a.cfg.App.Env, a.log)

	// Start pprof server on :6060
	// go func() {
	// 	if err := nethttp.ListenAndServe(":6060", nil); err != nil {
	// 		a.log.Errorf("Error starting pprof server: %s", err)
	// 	}
	// }()

	api.NewHandler(s.Rtr, a.c.useCase, a.log, a.c.GetOAuth2Authentik())

	go func() {
		err := metrics.Listen(":9100")
		if err != nil {
			a.log.Errorf("Error starting metrics server: %s", err)
		}
	}()

	err := s.Start(a.cfg.App.Port)
	if err != nil {
		return err
	}

	return nil
}
