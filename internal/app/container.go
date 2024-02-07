package app

import (
	"github.com/go-gorp/gorp"
	"github.com/redis/go-redis/v9"
	adClient "github.com/saygik/go-ad-client"
	"github.com/saygik/go-userinfo/internal/adapter/repository/ad"
	"github.com/saygik/go-userinfo/internal/adapter/repository/mssql"
	"github.com/saygik/go-userinfo/internal/adapter/repository/redisclient"
	"github.com/saygik/go-userinfo/internal/usecase"
)

type Container struct {
	mssql *gorp.DbMap
	glpi  *gorp.DbMap
	rc    *redis.Client
	ads   map[string]*adClient.ADClient
}

func NewAppContainer(mssqlConnect *gorp.DbMap, glpiConnect *gorp.DbMap, r *redis.Client, adclients map[string]*adClient.ADClient) *Container {
	c := &Container{
		mssql: mssqlConnect,
		glpi:  glpiConnect,
		rc:    r,
		ads:   adclients,
	}
	return c
}
func (c *Container) GetUseCase() *usecase.UseCase {
	return usecase.New(c.getMssqlRepository(), c.getRedisRepository(), c.getADRepository())
}

func (c *Container) getMssqlRepository() usecase.Repository {
	return mssql.NewRepository(c.mssql)
}

func (c *Container) getRedisRepository() usecase.Redis {
	return redisclient.New(c.rc)
}

func (c *Container) getADRepository() usecase.AD {
	return ad.New(c.ads)
}
