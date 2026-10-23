package main

import (
	"github.com/saygik/go-userinfo/cmd"
	"github.com/saygik/go-userinfo/internal/controller/http/api"
)

var (
	version   = "6.1.0" // переопределяется через -ldflags
	commit    = "none"
	buildDate = "23.10.2026"
)

func main() {
	api.SetBuildInfo(version, commit, buildDate)
	cmd.StartServer()
}
