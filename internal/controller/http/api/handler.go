package api

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	rtr *gin.Engine
	uc  UseCase
	log *logrus.Logger
}

type UseCase interface {
	Authenticate(entity.LoginForm) (bool, map[string]string, error)
}

func NewHandler(router *gin.Engine, uc UseCase, log *logrus.Logger) {
	h := &Handler{
		rtr: router,
		uc:  uc,
		log: log,
	}

	h.rtr.LoadHTMLGlob("./public/html/*")

	h.rtr.Static("/public", "./public")

	h.rtr.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginFindUserAPI": "v2.01",
			"goVersion":      runtime.Version(),
		})
	})

	h.rtr.NoRoute(h.NoRoute)
	authGroup := h.rtr.Group("/auth")

	authGroup.POST("/login", h.Login)
	authGroup.GET("/logout", h.Logout)
}

func (h *Handler) NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"Code": "404", "Message": "Not Found"})
	c.Abort()

}
