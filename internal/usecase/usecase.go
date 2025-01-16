package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	GetAppRoles() ([]entity.IdName, error)
	GetAppGroups() ([]entity.IdName, error)
	GetAppResources() ([]entity.IdName, error)
	GetDomainUsersIP(string) ([]entity.UserIPComputer, error)
	GetDomainUsersAvatars(string) ([]entity.IdNameAvatar, error)
	GetUserResourceAccess(string, string) (int, error)
	GetUserRoles(string) ([]entity.IdName, error)
	GetUserRole(string) entity.IdName
	SetUserRole(string, int) error
	GetUserGroups(string) ([]entity.IdName, error)
	AddUserGroup(string, int) error
	DelUserGroup(string, int) error
	GetUserAvatar(string) (string, error)
	GetCurrentUserResources(string) ([]entity.AppResource, error)
	GetUserSoftwares(string) ([]entity.IdName, error)
	GetUserSoftwaresByFio(string) ([]int64, error)
	GetUserActivity(string) ([]entity.UserActivity, error)
	AddOneUserSoftware(entity.SoftwareForm) error
	AddOneSoftwareUser(entity.SoftUser) error
	DelOneUserSoftware(string) error
	GetOrgCodes() ([]entity.OrgWithCodes, error)
	GetSoftwareUsers(int) ([]entity.SoftUser, error)
	SetUserAvatar(string, string) error
	SetUserIp(entity.UserActivityForm) (string, error)
	GetSchedule(int) (entity.Schedule, error)
	GetScheduleUsers(int, int) ([]entity.IdName, error)
	GetScheduleUserGroups(int, string) ([]entity.IdNameType, error)

	GetScheduleTasks(int) ([]entity.ScheduleTask, error)
	GetScheduleTasksNotifications(string) ([]entity.ScheduleTask, error)
	AddScheduleTask(entity.ScheduleTask) (entity.ScheduleTask, error)
	UpdateScheduleTask(entity.ScheduleTask) error
	UpdateScheduleTaskNotification(int) error
	UpdateScheduleTaskDisableMattermost(int) error
	DelScheduleTask(int) error
	GetOneDelegateMailBoxes(string) ([]entity.MailBoxDelegates, error)
	GetUserGlpiTrackingGroups(string) ([]entity.Id, error)
}
type Redis interface {
	ClearAllDomainsUsers()
	AddKeyFieldValue(string, string, []byte) error
	GetKeyFieldAll(string) (map[string]string, error)
	GetKeyFieldValue(string, string) (string, error)
	DelKeyField(string, string) error
	Delete(string) error
}
type AD interface {
	DomainList() []entity.DomainList
	GetDomainAdminsGLPI(string) int
	GetDomainUsers(string) ([]map[string]interface{}, error)
	GetDomainComputers(string) ([]map[string]interface{}, error)
	IsDomainExist(string) bool
	Authenticate(string, entity.LoginForm) (bool, map[string]string, error)
	GetGroupUsers(string, string) ([]map[string]interface{}, error)
	GetUser(string, string) (map[string]interface{}, error)
}

type GLPI interface {
	GetUserByName(string) (entity.GLPIUser, error)
	GetUserById(int) (entity.GLPIUser, error)
	GetUserProfiles(int) ([]entity.GLPIUserProfile, error)
	GetUserGroups(int) ([]entity.IdName, error)
	GetAllSoftwares() ([]entity.Software, error)
	GetSoftware(int) (entity.Software, error)
	GetTicketsNonClosed() ([]entity.GLPI_Ticket, error)
	GetTicketsNonClosedFromIniciator(int) ([]entity.GLPI_Ticket, error)
	GetSoftwaresAdmins() ([]entity.SoftwareAdmins, error)
	GetUsers() ([]entity.GLPIUserShort, error)
	GetGroupUsers(int) ([]entity.IdName, error)
	GetTicket(string) (entity.GLPI_Ticket, error)
	GetGLPITicketSolutionTemplates(string) ([]entity.GLPI_Ticket_Profile, error)
	GetTicketUsers(string) (users []entity.GLPIUser, err error)
	GetTicketGroupForCurrentUser(string, int) ([]entity.GLPIGroup, error)
	GetTicketWorks(string) ([]entity.GLPI_Work, error)
	GetTicketProblems(string) ([]entity.GLPI_Problem, error)
	GetProblemWorks(string) ([]entity.GLPI_Work, error)
	GetProblemTickets(string) ([]entity.GLPI_Otkaz, error)
	GetOtkazes(string, string) ([]entity.GLPI_Otkaz, error)
	GetProblems(string, string) ([]entity.GLPI_Problem, error)
	GetProblem(string) (entity.GLPI_Problem, error)
	GetProblemUsers(string) (users []entity.GLPIUser, err error)
	GetProblemGroups(string) ([]entity.IdName, error)
	GetProblemGroupForCurrentUser(string, int) ([]entity.GLPIGroup, error)
	GetStatTickets() ([]entity.GLPITicketsStats, error)
	GetStatFailures() ([]entity.GLPITicketsStats, error)
	GetStatPeriodRegionDayCounts(string, string, int) ([]entity.RegionsDayStats, error)
	GetStatTicketsDays(string, string) ([]entity.GLPITicketsStats, error)
	GetStatTop10Performers(string, string) ([]entity.GLPIStatsTop10, error)
	GetStatTop10Iniciators(string, string) ([]entity.GLPIStatsTop10, error)
	GetStatTop10Groups(string, string) ([]entity.GLPIStatsTop10, error)
	GetStatPeriodTicketsCounts(string, string) ([]entity.GLPIStatsCounts, error)
	GetStatPeriodRequestTypes(string, string) ([]entity.GLPIStatsTop10, error)
	GetStatRegions(string, string) ([]entity.GLPIRegionsStats, error)
	GetStatPeriodOrgTreemap(string, string) ([]entity.TreemapData, error)
	GetHRPTickets() ([]entity.GLPI_Ticket, error)
	GetHRPTicketsTest() ([]entity.GLPI_Ticket, error)
	SetHRPTicket(int) error
	GetUserApiTokenByName(string) (entity.IdName, error)
	GetTicketsInExecutionGroups(string) ([]entity.GLPI_Ticket, error)
	GetUserTrackingGroups(string) ([]entity.IdName, error)
	GetGroupMattermostChannel(int) (string, int, error)
}

type Mattermost interface {
	GetUsers() ([]entity.MattermostUserWithSessions, error)
	GetUserById(string) (*entity.MattermostUser, error)
	GetUsersWithSessions() (err error)
	IntegrationAllowedHosts() []string
	SendPostHRP(string, entity.HRPUser) (err error)
	SendPost(string, string, string, string, string, string, bool) (err error)
	SendPostHRPSoft(string, entity.HRPUser, entity.Software, int) (err error)
}
type GlpiApi interface {
	CreateTicket(entity.NewTicketInputForm) (int, error)
	CreateComment(entity.NewCommentInputForm) (int, error)
	CreateSolution(entity.NewCommentInputForm) (int, error)
	AddTicketUser(entity.GLPITicketUserInputForm) (int, error)
}

type UseCase struct {
	repo    Repository
	redis   Redis
	ad      AD
	glpi    GLPI
	matt    Mattermost
	glpiApi GlpiApi
	log     *logrus.Logger
}

func New(r Repository, redis Redis, adRepo AD, glpiRepo GLPI, matt Mattermost, glpiApi GlpiApi, log *logrus.Logger) *UseCase {
	return &UseCase{
		repo:    r,
		redis:   redis,
		ad:      adRepo,
		glpi:    glpiRepo,
		matt:    matt,
		glpiApi: glpiApi,
		log:     log,
	}
}
