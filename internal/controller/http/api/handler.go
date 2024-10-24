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
	GetAppRoles(string) ([]entity.IdName, error)                    // Все роли пользователей приложения
	SetUserRole(string, string, int) error                          // Установить роль пользователя
	GetAppGroups(string) ([]entity.IdName, error)                   // Все группы пользователей приложения
	DelUserGroup(string, string, int) error                         // Удалить группу приложения пользователя
	AddUserGroup(string, string, int) error                         // Добавить группу приложения пользователя
	GetAppResources(string) ([]entity.IdName, error)                // Все ресурсы приложения
	Authenticate(entity.LoginForm) (bool, map[string]string, error) // Аутентификация
	SetUserIp(entity.UserActivityForm) (string, error)
	GetADUsers(string) ([]map[string]interface{}, error)                                // Пользователи домена
	GetADUsersPublicInfo(string) ([]map[string]interface{}, error)                      // Пользователи домена соклащённая информация
	GetADComputers(string) ([]map[string]interface{}, error)                            // Компьютеры домена
	GetUser(string, string) (map[string]interface{}, error)                             // Свойства пользователя домена
	GetCurrentUser(string, string) (map[string]interface{}, error)                      // Свойства залогиненного пользователя домена
	GetUserADPropertys(string, string) (map[string]interface{}, error)                  // Разрешённые сппециалисту свойства пользователя домена
	GetCurrentUserResources(string) ([]entity.AppResource, error)                       // Разрешённые ресурсы
	GetGlpiUser(string) (entity.GLPIUser, error)                                        // Пользователь GLPI
	GetGlpiUserForTechnical(string, string) (*entity.GLPIUser, error)                   // Пользователь GLPI для технического специалиста
	GetAdCounts() (int, int, error)                                                     // Основная статистика доменов
	GetDomainList(string) []entity.DomainList                                           // Список доменов
	GetADGroupUsers(string, string) ([]map[string]interface{}, error)                   //
	GetUserSoftwares(string) ([]entity.Software, error)                                 // Список систем пользователя
	GetSoftwares() ([]entity.Software, error)                                           // Список всех систем
	GetSoftware(string) (entity.Software, error)                                        // Одна система
	GetSoftwareUsers(string) ([]map[string]interface{}, error)                          // Список пользователей одной системы
	AddUserSoftware(entity.SoftwareForm) error                                          // Добавление системы пользователя
	AddOneSoftwareUser(string, entity.SoftUser) (map[string]interface{}, error)         // Добавление пользователя в систему
	DelUserSoftware(entity.SoftwareForm) error                                          // Удаление системы пользователя
	GetUserADActivity(string, string) ([]entity.UserActivity, error)                    // Активность доменов
	SetUserAvatar(string, string, string) error                                         // Установить пользователю аватар
	GetOrgCodes() ([]entity.OrgWithCodes, error)                                        // Коды организаций
	GetMattermostUsers() ([]entity.MattermostUser, error)                               // Все пользователи Mattermost
	GetGLPITicketsNonClosed(string) ([]entity.GLPI_Ticket, error)                       // Все нерешённые заявки GLPI
	GetGLPIUsers() ([]entity.GLPIUserShort, error)                                      // Все пользователи GLPI
	GetGLPITicket(string, string) (entity.GLPI_Ticket, error)                           // Одна заявка GLPI
	GetGLPITicketSolutionTemplates(string) ([]entity.GLPI_Ticket_Profile, error)        // Шаблоны решений заявки
	GetGLPIProblem(string, string) (entity.GLPI_Problem, error)                         // Одна проблема GLPI
	GetGLPIOtkazes(string, string) ([]entity.GLPI_Otkaz, error)                         // Отказы GLPI за период
	GetGLPIProblems(string, string) ([]entity.GLPI_Problem, error)                      // Проблемы GLPI за период
	GetStatTickets() ([]entity.GLPITicketsStats, error)                                 //
	GetStatFailures() ([]entity.GLPITicketsStats, error)                                //
	GetStatPeriodRegionDayCounts(string, string, int) ([]entity.RegionsDayStats, error) //
	GetStatTicketsDays(string, string) ([]entity.GLPITicketsStats, error)               //
	GetStatTop10Performers(string, string) ([]entity.GLPIStatsTop10, error)             //
	GetStatTop10Iniciators(string, string) ([]entity.GLPIStatsTop10, error)             //
	GetStatTop10Groups(string, string) ([]entity.GLPIStatsTop10, error)                 //
	GetStatPeriodTicketsCounts(string, string) ([]entity.GLPIStatsCounts, error)        //
	GetStatPeriodRequestTypes(string, string) ([]entity.GLPIStatsTop10, error)          //
	GetStatRegions(string, string) ([]entity.GLPIRegionsStats, error)                   //
	GetStatPeriodOrgTreemap(string, string) ([]entity.TreemapData, error)               //
	GetSchedule(string) (entity.Schedule, error)                                        // Один календарь
	GetScheduleTasks(string) ([]entity.ScheduleTask, error)                             //
	AddScheduleTask(entity.ScheduleTask) (entity.ScheduleTask, error)                   //
	DelScheduleTask(string) error                                                       //
	UpdateScheduleTask(entity.ScheduleTask) error                                       //
	AddTicketSolution(entity.NewCommentForm) error                                      // GLPI. Добавление  решения
	AddTicketComment(entity.NewCommentForm) error                                       // GLPI. Добавление  комментария
	AddTicket(entity.NewTicketForm) (int, error)                                        // GLPI. Добавление  заявки
	AddTicketUser(entity.GLPITicketUserForm) error                                      //
	GetTicketsInExecutionGroups(string) ([]entity.GLPI_Ticket, error)                   // Незакрытые заявки в группах слежения пользователя
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
	h.NewSoftwareRouterGroup()
	h.NewManualRouterGroup()
	h.NewMattermostRouterGroup()
	h.NewScheduleRouterGroup()
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
