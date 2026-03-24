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
	oAuth2Authentik OAuth2
	allowedApi      []string
}

type UseCase interface {
	ComputerLocalAdminsAudit(string, []string, bool) error
	ComputerLocalAdminsGet(bool) ([]entity.LocalAdmins, error)
	UpdateComputerLocalAdmins(entity.Permissions, string, string, string) error
	FillRedisCaсheFromAD() error
	LoadUserPermissions(string) (entity.Permissions, error)                                                      //LoadUserPermissions
	IUTMGetWlist() ([]string, error)                                                                             //IUTMGetWlist
	IsAppInitialized() bool                                                                                      // Приложение инициализировано
	ADUserLocked(string) bool                                                                                    // Пользователь заблокирован
	GetAppRoles() ([]entity.IdNameDescription, error)                                                            // Все роли пользователей приложения
	GetSections(string) ([]entity.IdNameDescription, error)                                                      // Все роли пользователей приложения
	GetDomainAccess(string) ([]entity.DomainAccess, error)                                                       // Все роли пользователей приложения
	GetUserRoles(string) ([]entity.IdNameDescription, error)                                                     // Все роли пользователей приложения
	GetAppSections() ([]entity.IdNameDescription, error)                                                         // Все разделы  приложения
	SetUserRole(entity.Permissions, string, int) error                                                           // Установить роль пользователя
	DelUserRole(entity.Permissions, string, int) error                                                           // Удалить роль приложения пользователя
	DelUserDomainRole(entity.Permissions, string, string) error                                                  // Удалить роль приложения пользователя
	DelUserSection(entity.Permissions, string, string) error                                                     // Удалить раздел приложения пользователя
	AddUserRole(entity.Permissions, string, int) (*entity.IdNameDescription, error)                              // Добавить роль приложения пользователя
	AddUserDomainRole(entity.Permissions, string, string, string) (*entity.DomainAccess, string, error)          // Добавить роль приложения пользователя
	AddUserSection(entity.Permissions, string, string) (*entity.IdNameDescription, error)                        // Добавить раздел приложения пользователя
	GetAppResources(string) ([]entity.IdName, error)                                                             // Все ресурсы приложения
	Authenticate(entity.LoginForm) (bool, map[string]string, error)                                              // Аутентификация
	SetUserIp(entity.UserActivityForm) (string, error)                                                           // Установить ip пользователя
	GetADUsers(entity.Permissions) ([]map[string]interface{}, error)                                             // Пользователи домена
	GetADUsersPublicInfo(entity.Permissions) ([]map[string]interface{}, error)                                   // Пользователи домена соклащённая информация
	GetADComputers(entity.Permissions) ([]map[string]interface{}, error)                                         // Компьютеры домена
	GetADComputersVersions(string, entity.Permissions) ([]entity.ComputerVersionCount, error)                    // Кол-во компьютеров домена по версиям ОС
	GetADComputersOSFamily(string, entity.Permissions) ([]entity.ComputerFamilyCount, error)                     // Кол-во компьютеров домена по семействам ОС
	GetADLastComputers(string, entity.Permissions) ([]entity.DomainComputer, error)                              // Компьютеры домена с последними логинами
	GetUser(string) (map[string]interface{}, error)                                                              // Свойства пользователя домена
	UserExist(string) error                                                                                      // Существует ли пользователь в доменах
	GetCurrentUser(entity.Permissions) (map[string]interface{}, error)                                           // Свойства залогиненного пользователя домена
	GetUserADPropertys(string, entity.Permissions) (map[string]interface{}, error)                               // Разрешённые сппециалисту свойства пользователя домена
	GetUserADPropertysSimple(string) (*entity.SimpleUser, error)                                                 // Упрощенные  свойства пользователя домена
	GetCurrentUserResources(string) ([]entity.AppResource, error)                                                // Разрешённые ресурсы
	GetGlpiUser(string) (entity.GLPIUser, error)                                                                 // Пользователь GLPI
	GetGlpiUserForTechnical(string, string) (*entity.GLPIUser, error)                                            // Пользователь GLPI для технического специалиста
	GetAdCounts() (int, int, error)                                                                              // Основная статистика доменов (все домены)
	GetAdCountsDomain(string) (int, int, error)                                                                  // Статистика по одному домену
	GetDomainList(entity.Permissions) []entity.DomainList                                                        // Список доменов
	DomainList() []entity.DomainList                                                                             // Список всех доменов
	GetADGroupUsers(string, string) ([]map[string]interface{}, error)                                            // Пользователи группы
	UserInDomainGroup2(string, string, string) error                                                             // Проверить, находится ли пользователь в группе домена
	ADUserAddGroup(string, string) error                                                                         // Добавить пользователя в группу
	ADUserDelGroup(string, string) error                                                                         // Удалить пользователя из группы
	SwitchUserGroupInternet(string, string, string, bool, int) error                                             // Изменить группу интернета пользователя (с поддержкой временного режима)
	GetTemporaryGroupChange(string) (*entity.TemporaryGroupChange, error)                                        // Получить информацию о временном изменении группы пользователя
	DeleteTemporaryGroupChange(string) error                                                                     // Удалить временное изменение группы и восстановить предыдущую группу
	GetUserSoftwares(string) ([]entity.Software, error)                                                          // Список систем пользователя
	GetSoftwares() ([]entity.Software, error)                                                                    // Список всех систем
	GetSoftwaresUsers() ([]entity.SoftUser, error)                                                               // Список всех систем пользователя
	GetSoftware(string) (entity.Software, error)                                                                 // Одна система
	GetSoftwareUsers(string) ([]map[string]interface{}, error)                                                   // Список пользователей одной системы
	AddUserSoftware(entity.SoftwareForm) error                                                                   // Добавление системы пользователя
	AddOneSoftwareUser(string, entity.SoftUser, string) (map[string]interface{}, error)                          // Добавление пользователя в систему
	UpdateOneSoftwareUser(entity.SoftUser, string) (entity.SoftUser, error)                                      // Изменение пользователя в систему
	DelUserSoftware(string) error                                                                                // Удаление системы пользователя
	GetUserADActivity(string, entity.Permissions) ([]entity.UserActivity, error)                                 // Активность доменов
	GetUserMailboxPermissions(string, entity.Permissions) ([]entity.MailBoxDelegates, error)                     // Получение делегированных почтовых ящиков пользователя
	SetUserAvatar(string, string, string) error                                                                  // Установить пользователю аватар
	GetOrgCodes() ([]entity.OrgWithCodes, error)                                                                 // Коды организаций
	GetMattermostUsers() ([]entity.MattermostUserWithSessions, error)                                            // Все пользователи Mattermost
	AddGLPI_HRPTicketCommentFromMattermost(entity.MattermostInteractiveMessageRequestForm) (string, error)       // Добавить комментарий из Mattermost в GLPI заявку кадровичка
	DisableSheduleTaskNotificationFromMattermost(entity.MattermostInteractiveMessageRequestForm) (string, error) // Выключить оповещение задачи календаря из Mattermost
	MattermostIntegrationAllowed(string) bool                                                                    // Запрос с этого ip разрешен для интеграции
	GetGLPITicketsNonClosed(string) ([]entity.GLPI_Ticket, error)                                                // Все нерешённые заявки GLPI
	GetGLPIUsers() ([]entity.GLPIUserShort, error)                                                               // Все пользователи GLPI
	GetGLPITicket(string, string) (entity.GLPI_Ticket, error)                                                    // Одна заявка GLPI
	GetGLPITicketReport(string, string) (*entity.GLPI_Ticket_Report, error)                                      // Одна заявка GLPI для отчетов
	GetGLPITicketSimple(string) (entity.GLPI_Ticket, error)                                                      // Одна заявка GLPI
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
	GetComputerRMS(string) ([]entity.ComputerRMS, error)                                                         //Списсок компьютеров с RMS
	ReplyPost(string, string, string) error                                                                      // Mattermost. Ответить на пост
}

type OAuth2 interface {
	AuthCodeURL(string) string
	LogOutURL() string
	Exchange(string) (*entity.Token, *entity.UserInfo, error)
	IntrospectOAuth2Token(string, bool) (*entity.UserInfo, error)
	ExchangeRefreshToAccessToken(string) (*entity.Token, error)
	GetRedirectUrl() string
	Refresh(string) (entity.Token, error)
}

func NewHandler(router *gin.Engine, uc UseCase, log *logrus.Logger, oAuth2Authentik OAuth2, allowedApi []string) {
	h := &Handler{
		rtr:        router,
		uc:         uc,
		log:        log,
		allowedApi: allowedApi,

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
	h.NewOAuth2RouterGroup()

	h.NewADRouterGroup()
	h.NewImgRouterGroup()
	h.NewAppRouterGroup()
	h.NewGlpiRouterGroup()
	h.NewSoftwareRouterGroup()
	h.NewManualRouterGroup()
	h.NewMattermostRouterGroup()
	h.NewScheduleRouterGroup()
	h.NewMattermostCommandsRouterGroup()
	h.NewIUTMRouterGroup()
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
