package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/saygik/go-userinfo/internal/app"
)

func StartServer() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	a, err := app.New(ctx)
	if err != nil {
		log.Fatalf("Config error: %s", err)
		os.Exit(1)
	}
	err = a.StartHTTPServer()
	if err != nil {
		log.Fatalf("Error starting WEB-server: %s", err)
		os.Exit(2)
	}
	<-ctx.Done()
}
