package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

type Server struct {
	Rtr *gin.Engine
}

func NewServer(env string, log *logrus.Logger) *Server {
	server := &Server{}
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	server.Rtr = gin.New()
	server.Rtr.Use(ginlogrus.Logger(log), gin.Recovery())
	server.registerMiddlewares()

	return server
}

func (s *Server) Start(port string) error {

	err := http.ListenAndServe(":"+port, s.Rtr)
	if err != nil {
		return err
	}
	return nil
}
