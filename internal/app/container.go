package app

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-gorp/gorp"
	"github.com/redis/go-redis/v9"
	adClient "github.com/saygik/go-ad-client"
	"github.com/saygik/go-userinfo/internal/adapter/repository/ad"
	"github.com/saygik/go-userinfo/internal/adapter/repository/glpi"
	"github.com/saygik/go-userinfo/internal/adapter/repository/mssql"
	"github.com/saygik/go-userinfo/internal/adapter/repository/redisclient"
	"github.com/saygik/go-userinfo/internal/adapter/web/glpiapi"
	"github.com/saygik/go-userinfo/internal/adapter/web/mattermost"
	"github.com/saygik/go-userinfo/internal/auth/hydra"
	"github.com/saygik/go-userinfo/internal/auth/oauth2"
	"github.com/saygik/go-userinfo/internal/usecase"
)

type Container struct {
	mssql   *gorp.DbMap
	glpi    *gorp.DbMap
	rc      *redis.Client
	ads     map[string]*adClient.ADClient
	matt    *MattClient
	glpiApi *GLPIApiClient
	hydra   *IDPClient
	oAuth2  *OAuth2Client
}

func NewAppContainer(
	mssqlConnect *gorp.DbMap,
	glpiConnect *gorp.DbMap,
	r *redis.Client,
	adclients map[string]*adClient.ADClient,
	matt *MattClient,
	glpiApi *GLPIApiClient,
	hydra *IDPClient,
	oAuth2 *OAuth2Client,
) *Container {
	c := &Container{
		mssql:   mssqlConnect,
		glpi:    glpiConnect,
		rc:      r,
		ads:     adclients,
		matt:    matt,
		glpiApi: glpiApi,
		hydra:   hydra,
		oAuth2:  oAuth2,
	}
	return c
}
func (c *Container) GetUseCase() *usecase.UseCase {
	return usecase.New(c.getMssqlRepository(), c.getRedisRepository(), c.getADRepository(), c.getGlpiRepository(), c.getMattermostRepository(), c.getGlpiApiRepository())
}

func (c *Container) getMssqlRepository() *mssql.Repository {
	return mssql.NewRepository(c.mssql)
}
func (c *Container) getGlpiRepository() *glpi.Repository {
	return glpi.NewRepository(c.glpi)
}

func (c *Container) getRedisRepository() *redisclient.Repository {
	return redisclient.New(c.rc)
}

func (c *Container) getMattermostRepository() *mattermost.Repository {
	return mattermost.New(c.matt.url, c.matt.token)
}
func (c *Container) getGlpiApiRepository() *glpiapi.Repository {
	return glpiapi.New(c.glpiApi.url, c.glpiApi.token, c.glpiApi.usertoken)
}

func (c *Container) GetHydra() *hydra.Hydra {
	return hydra.New(c.hydra.client, c.hydra.scopes)
}
func (c *Container) GetOAuth2() *oauth2.OAuth2 {
	return oauth2.New(c.oAuth2.oAuth2Config, c.oAuth2.oidcProvider)
}

func (c *Container) getADRepository() *ad.Repository {
	return ad.New(c.ads)
}
