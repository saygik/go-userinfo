package usecase

import (
	"github.com/saygik/go-userinfo/internal/entity"
)

type Repository interface {
	GetAllDomainsUsers()
	GetDomainUsersIP(string) ([]entity.UserIPComputer, error)
	GetDomainUsersAvatars(string) ([]entity.IdNameAvatar, error)
}
type Redis interface {
	ClearAllDomainsUsers()
	AddKeyFieldValue(string, string, []byte) error
}
type AD interface {
	DomainList() []string
	GetDomainUsers(string) ([]map[string]interface{}, error)
	GetDomainComputers(string) ([]map[string]interface{}, error)
	IsDomainExist(string) bool
	Authenticate(string, entity.LoginForm) (bool, map[string]string, error)
}

type UseCase struct {
	repo  Repository
	redis Redis
	ad    AD
}

func New(r Repository, redis Redis, adRepo AD) *UseCase {
	return &UseCase{
		repo:  r,
		redis: redis,
		ad:    adRepo,
	}
}
