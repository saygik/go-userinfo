package cmd

import (
	"log"
	"os"

	"github.com/saygik/go-userinfo/internal/app"
)

func StartServer() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
		os.Exit(1)
	}
	err = a.StartHTTPServer()
	if err != nil {
		log.Fatalf("Error starting WEB-server: %s", err)
		os.Exit(2)
	}
}
