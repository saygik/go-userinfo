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
	rg  *gin.RouterGroup
	uc  UseCase
	log *logrus.Logger
	jwt JWT
}

type UseCase interface {
	Authenticate(entity.LoginForm) (bool, map[string]string, error)    // Аутентификация
	GetADUsers(string) ([]map[string]interface{}, error)               // Пользователи домена
	GetADComputers(string) ([]map[string]interface{}, error)           // Компьютеры домена
	GetUser(string) (map[string]interface{}, error)                    // Свойства пользователя домена
	GetUserADPropertys(string, string) (map[string]interface{}, error) // Разрешённые сппециалисту свойства пользователя домена
	GetCurrentUserResources(string) ([]entity.AppResource, error)      // Разрешённые ресурсы
	GetGlpiUser(string) (entity.GLPIUser, error)                       // Пользователь GLPI
	GetGlpiUserForTechnical(string, string) (entity.GLPIUser, error)   // Пользователь GLPI для технического специалиста
	GetAdCounts() (int, int, error)                                    // Основная статистика доменов
	GetDomainList(string) []entity.DomainList                          // Список доменов
	GetUserSoftwares(string) ([]entity.Software, error)
}

type JWT interface {
	Login(string) (entity.Token, error)
	DeleteAuth(string) error
	ExtractTokenMetadata(r *http.Request) (*entity.AccessDetails, error)
	FetchAuth(*entity.AccessDetails) (string, error)
	RefreshToken(string) (map[string]string, error)
}

func NewHandler(router *gin.Engine, uc UseCase, log *logrus.Logger, jw JWT) {
	h := &Handler{
		rtr: router,
		uc:  uc,
		log: log,
		jwt: jw,
	}

	h.rtr.LoadHTMLGlob("./public/html/*")

	h.rtr.Static("/public", "./public")

	h.rtr.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginFindUserAPI": "v2.01",
			"goVersion":      runtime.Version(),
		})
	})

	h.rg = h.rtr.Group("/api")
	h.NewAuthRouterGroup()
	h.NewADRouterGroup()
	h.NewAppRouterGroup()
	h.NewGlpiRouterGroup()
	h.rtr.NoRoute(h.NoRoute)
}

func (h *Handler) NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"Code": "404", "Message": "Not Found"})
	c.Abort()

}

// TokenAuthMiddleware ...
// JWT Authentication middleware attached to each request that needs to be authenitcated to validate the access_token in the header
func (h *Handler) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.TokenValid(c)
		c.Next()
	}
}
