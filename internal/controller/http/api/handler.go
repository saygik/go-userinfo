package api

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/entity"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	rtr             *gin.Engine
	rg              *gin.RouterGroup
	uc              UseCase
	log             *logrus.Logger
	hydra           Hydra
	oAuth2          OAuth2
	oAuth2Authentik OAuth2
}

type UseCase interface {
	ADUserLocked(string) bool
	GetAppRoles(string) ([]entity.IdName, error)                    // Все роли пользователей приложения
	SetUserRole(string, string, int) error                          // Установить роль пользователя
	GetAppGroups(string) ([]entity.IdName, error)                   // Все группы пользователей приложения
	DelUserGroup(string, string, int) error                         // Удалить группу приложения пользователя
	AddUserGroup(string, string, int) error                         // Добавить группу приложения пользователя
	DelUserRole(string, string, int) error                          // Удалить роль приложения пользователя
	AddUserRole(string, string, int) error                          // Добавить роль приложения пользователя
	GetAppResources(string) ([]entity.IdName, error)                // Все ресурсы приложения
	Authenticate(entity.LoginForm) (bool, map[string]string, error) // Аутентификация
	SetUserIp(entity.UserActivityForm) (string, error)
	GetADUsers(string) ([]map[string]interface{}, error)                                                         // Пользователи домена
	GetADUsersPublicInfo(string) ([]map[string]interface{}, error)                                               // Пользователи домена соклащённая информация
	GetADComputers(string) ([]map[string]interface{}, error)                                                     // Компьютеры домена
	GetUser(string, string) (map[string]interface{}, error)                                                      // Свойства пользователя домена
	UserExist(string) error                                                                                      // Существует ли пользователь в доменах
	GetCurrentUser(string, string) (map[string]interface{}, error)                                               // Свойства залогиненного пользователя домена
	GetUserADPropertys(string, string) (map[string]interface{}, error)                                           // Разрешённые сппециалисту свойства пользователя домена
	GetCurrentUserResources(string) ([]entity.AppResource, error)                                                // Разрешённые ресурсы
	GetGlpiUser(string) (entity.GLPIUser, error)                                                                 // Пользователь GLPI
	GetGlpiUserForTechnical(string, string) (*entity.GLPIUser, error)                                            // Пользователь GLPI для технического специалиста
	GetAdCounts() (int, int, error)                                                                              // Основная статистика доменов
	GetDomainList(string) []entity.DomainList                                                                    // Список доменов
	GetADGroupUsers(string, string) ([]map[string]interface{}, error)                                            //
	GetUserSoftwares(string) ([]entity.Software, error)                                                          // Список систем пользователя
	GetSoftwares() ([]entity.Software, error)                                                                    // Список всех систем
	GetSoftwaresUsers() ([]entity.SoftUser, error)                                                               // Список всех систем пользователя
	GetSoftware(string) (entity.Software, error)                                                                 // Одна система
	GetSoftwareUsers(string) ([]map[string]interface{}, error)                                                   // Список пользователей одной системы
	AddUserSoftware(entity.SoftwareForm) error                                                                   // Добавление системы пользователя
	AddOneSoftwareUser(string, entity.SoftUser, string) (map[string]interface{}, error)                          // Добавление пользователя в систему
	UpdateOneSoftwareUser(entity.SoftUser, string) (entity.SoftUser, error)                                      // Изменение пользователя в систему
	DelUserSoftware(string) error                                                                                // Удаление системы пользователя
	GetUserADActivity(string, string) ([]entity.UserActivity, error)                                             // Активность доменов
	GetUserMailboxPermissions(string, string) ([]entity.MailBoxDelegates, error)                                 // Получение делегированных почтовых ящиков пользователя
	SetUserAvatar(string, string, string) error                                                                  // Установить пользователю аватар
	GetOrgCodes() ([]entity.OrgWithCodes, error)                                                                 // Коды организаций
	GetMattermostUsers() ([]entity.MattermostUserWithSessions, error)                                            // Все пользователи Mattermost
	AddGLPI_HRPTicketCommentFromMattermost(entity.MattermostInteractiveMessageRequestForm) (string, error)       // Добавить комментарий из Mattermost в GLPI заявку кадровичка
	DisableSheduleTaskNotificationFromMattermost(entity.MattermostInteractiveMessageRequestForm) (string, error) // Выключить оповещение задачи календаря из Mattermost
	MattermostIntegrationAllowed(string) bool                                                                    // Запрос с этого ip разрешен для интеграции
	GetGLPITicketsNonClosed(string) ([]entity.GLPI_Ticket, error)                                                // Все нерешённые заявки GLPI
	GetGLPIUsers() ([]entity.GLPIUserShort, error)                                                               // Все пользователи GLPI
	GetGLPITicket(string, string) (entity.GLPI_Ticket, error)                                                    // Одна заявка GLPI
	GetGLPITicketSolutionTemplates(string) ([]entity.GLPI_Ticket_Profile, error)                                 // Шаблоны решений заявки
	GetGLPIProblem(string, string) (entity.GLPI_Problem, error)                                                  // Одна проблема GLPI
	GetGLPIOtkazes(string, string) ([]entity.GLPI_Otkaz, error)                                                  // Отказы GLPI за период
	GetGLPIProblems(string, string) ([]entity.GLPI_Problem, error)                                               // Проблемы GLPI за период
	GetStatTickets() ([]entity.GLPITicketsStats, error)                                                          //
	GetStatFailures() ([]entity.GLPITicketsStats, error)                                                         //
	GetStatPeriodRegionDayCounts(string, string, int) ([]entity.RegionsDayStats, error)                          //
	GetStatTicketsDays(string, string) ([]entity.GLPITicketsStats, error)                                        //
	GetStatTop10Performers(string, string) ([]entity.GLPIStatsTop10, error)                                      //
	GetStatTop10Iniciators(string, string) ([]entity.GLPIStatsTop10, error)                                      //
	GetStatTop10Groups(string, string) ([]entity.GLPIStatsTop10, error)                                          //
	GetStatPeriodTicketsCounts(string, string) ([]entity.GLPIStatsCounts, error)                                 //
	GetStatPeriodRequestTypes(string, string) ([]entity.GLPIStatsTop10, error)                                   //
	GetStatRegions(string, string) ([]entity.GLPIRegionsStats, error)                                            //
	GetStatPeriodOrgTreemap(string, string) ([]entity.TreemapData, error)                                        //
	GetSchedule(int, string) (entity.Schedule, error)                                                            // Один календарь
	GetAllSchedules(string) ([]entity.IdNameType, error)                                                         // Все календарь
	GetScheduleTasks(int) ([]entity.ScheduleTaskCalendar, error)                                                 //
	AddScheduleTask(entity.ScheduleTask) (entity.ScheduleTaskCalendar, error)                                    //
	DelScheduleTask(int) error                                                                                   //
	UpdateScheduleTask(entity.ScheduleTask) error                                                                //
	AddTicketSolution(entity.NewCommentForm) error                                                               // GLPI. Добавление  решения
	AddTicketComment(entity.NewCommentForm) error                                                                // GLPI. Добавление  комментария
	AddTicket(entity.NewTicketForm) (int, error)                                                                 // GLPI. Добавление  заявки
	AddTicketUser(entity.GLPITicketUserForm) error                                                               //
	GetTicketsInExecutionGroups(string) ([]entity.GLPI_Ticket, error)                                            // Незакрытые заявки в группах слежения пользователя
	UserInGropScopes(string, []string, []entity.IDPScope, *entity.OAuth2Client) ([]string, bool, error)          // У пользователя есть права на требуемый scope
}

type Hydra interface {
	CheckHydra() bool
	GetOAuth2LoginRequest(string) (*entity.OAuth2LoginRequest, error)
	AcceptOAuth2LoginRequest(string, string) (string, error)
	GetOAuth2ConsentRequest(string) (*entity.OAuth2ConsentRequest, error)
	AcceptNewOAuth2LoginRequest(string, string, bool) (string, error)
	AcceptOAuth2ConsentRequest(*entity.OAuth2ConsentRequest, map[string]interface{}) (string, error)
	IntrospectOAuth2Token(string) (*entity.IntrospectedOAuth2Token, error)
	AcceptOAuth2LogoutRequest(string) (string, error)
	GetOAuth2Client(string) (*entity.OAuth2Client, error)
	LogoutURL() string

	GetScopes() []entity.IDPScope
}

type OAuth2 interface {
	AuthCodeURL(string) string
	LogOutURL() string
	Exchange(string) (*entity.Token, *entity.UserInfo, error)
	IntrospectOAuth2Token(string) (*entity.UserInfo, error)
	GetRedirectUrl() string
	Refresh(string) (entity.Token, error)
}

func NewHandler(router *gin.Engine, uc UseCase, log *logrus.Logger, hydra Hydra, oAuth2 OAuth2, oAuth2Authentik OAuth2) {
	h := &Handler{
		rtr:             router,
		uc:              uc,
		log:             log,
		hydra:           hydra,
		oAuth2:          oAuth2,
		oAuth2Authentik: oAuth2Authentik,
	}

	h.rtr.Static("/public", "./public")
	router.Static("/css", "./public/css")
	router.Static("/js", "./public/js")
	h.rtr.LoadHTMLGlob("./public/html/*")

	h.rtr.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginFindUserAPI": "v2.01",
			"goVersion":      runtime.Version(),
		})
	})
	h.rtr.GET("/f", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"goVersion":      runtime.Version(),
			"LoginChallenge": "LoginChallenge",
		})
	})
	h.rg = h.rtr.Group("/api")
	h.NewOAuthRouterGroup()
	h.NewOAuth2RouterGroup()
	h.NewHydraIDPRouterGroup()
	h.NewADRouterGroup()
	h.NewAppRouterGroup()
	h.NewGlpiRouterGroup()
	h.NewSoftwareRouterGroup()
	h.NewManualRouterGroup()
	h.NewMattermostRouterGroup()
	h.NewScheduleRouterGroup()
	h.NewMattermostCommandsRouterGroup()
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

func (h *Handler) UserFromTokenTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.UserFromToken(c)
		c.Next()
	}
}
