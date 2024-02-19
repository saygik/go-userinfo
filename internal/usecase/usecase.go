package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

type Repository interface {
	GetDomainUsersIP(string) ([]entity.UserIPComputer, error)
	GetDomainUsersAvatars(string) ([]entity.IdNameAvatar, error)
	GetUserResourceAccess(string, string) (int, error)
	GetUserRoles(string) ([]entity.IdName, error)
	GetUserGroups(string) ([]entity.IdName, error)
	GetUserAvatar(string) (string, error)
	GetCurrentUserResources(string) ([]entity.AppResource, error)
}
type Redis interface {
	ClearAllDomainsUsers()
	AddKeyFieldValue(string, string, []byte) error
	GetKeyFieldAll(string) (map[string]string, error)
	GetKeyFieldValue(string, string) (string, error)
	DelKeyField(string, string) error
}
type AD interface {
	DomainList() []entity.DomainList
	GetDomainUsers(string) ([]map[string]interface{}, error)
	GetDomainComputers(string) ([]map[string]interface{}, error)
	IsDomainExist(string) bool
	Authenticate(string, entity.LoginForm) (bool, map[string]string, error)
}

type GLPI interface {
	GetUserByName(string) (entity.GLPIUser, error)
	GetUserProfiles(int64) ([]entity.GLPIUserProfile, error)
	GetUserGroups(int64) ([]entity.IdName, error)
	GetAllSoftwares() ([]entity.Software, error)
	GetSoftwaresAdmins() ([]entity.SoftwareAdmins, error)
}

type UseCase struct {
	repo  Repository
	redis Redis
	ad    AD
	glpi  GLPI
}

func New(r Repository, redis Redis, adRepo AD, glpiRepo GLPI) *UseCase {
	return &UseCase{
		repo:  r,
		redis: redis,
		ad:    adRepo,
		glpi:  glpiRepo,
	}
}
