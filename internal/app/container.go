package app

import (
	"context"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-gorp/gorp"
	"github.com/redis/go-redis/v9"
	"github.com/saygik/go-userinfo/internal/adapter/repository/ad"
	"github.com/saygik/go-userinfo/internal/adapter/repository/glpi"
	"github.com/saygik/go-userinfo/internal/adapter/repository/mssql"
	"github.com/saygik/go-userinfo/internal/adapter/repository/redisclient"
	"github.com/saygik/go-userinfo/internal/adapter/web/glpiapi"
	"github.com/saygik/go-userinfo/internal/adapter/web/mattermost"
	"github.com/saygik/go-userinfo/internal/adapter/web/webclient"
	"github.com/saygik/go-userinfo/internal/auth/hydra"
	"github.com/saygik/go-userinfo/internal/auth/oauth2"
	"github.com/saygik/go-userinfo/internal/auth/oauth2authentik"
	"github.com/saygik/go-userinfo/internal/config"
	"github.com/saygik/go-userinfo/internal/usecase"
	adClient "github.com/saygik/go-userinfo/pkg/ad-client"
	"github.com/sirupsen/logrus"
)

type Container struct {
	ctx             context.Context
	mssql           *gorp.DbMap
	glpi            *gorp.DbMap
	rc              *redis.Client
	ads             map[string]*adClient.ADClient
	adconfigs       map[string]*config.ADConfig
	matt            *MattClient
	glpiApi         *GLPIApiClient
	hydra           *IDPClient
	oAuth2          *OAuth2Client
	oAuth2Authentik *OAuth2Client
	webClientUrl    string
	log             *logrus.Logger
	// cached singletons
	useCase *usecase.UseCase
}

func NewAppContainer(
	ctx context.Context,
	mssqlConnect *gorp.DbMap,
	glpiConnect *gorp.DbMap,
	r *redis.Client,
	adclients map[string]*adClient.ADClient,
	adconfigs map[string]*config.ADConfig,
	matt *MattClient,
	glpiApi *GLPIApiClient,
	hydra *IDPClient,
	oAuth2 *OAuth2Client,
	oAuth2Authentik *OAuth2Client,
	webClientUrl string,
	log *logrus.Logger,
) *Container {
	c := &Container{
		ctx:             ctx,
		mssql:           mssqlConnect,
		glpi:            glpiConnect,
		rc:              r,
		ads:             adclients,
		adconfigs:       adconfigs,
		matt:            matt,
		glpiApi:         glpiApi,
		hydra:           hydra,
		oAuth2:          oAuth2,
		oAuth2Authentik: oAuth2Authentik,
		webClientUrl:    webClientUrl,
		log:             log,
	}
	uc := usecase.New(
		c.getMssqlRepository(),
		c.getRedisRepository(),
		c.getADRepository(),
		c.getGlpiRepository(),
		c.getMattermostRepository(),
		c.getGlpiApiRepository(),
		c.getWebClientRepository(),
		c.log,
	)
	c.useCase = uc
	return c
}

// func (c *Container) GetUseCase() *usecase.UseCase {
// 	if c.useCase == nil {
// 		c.useCase = usecase.New(
// 			c.getMssqlRepository(),
// 			c.getRedisRepository(),
// 			c.getADRepository(),
// 			c.getGlpiRepository(),
// 			c.getMattermostRepository(),
// 			c.getGlpiApiRepository(),
// 			c.getWebClientRepository(),
// 			c.log,
// 		)
// 	}
// 	return c.useCase
// }

func (c *Container) getWebClientRepository() *webclient.Repository {
	return webclient.New(c.ctx, c.webClientUrl, c.log)
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
	return mattermost.New(c.matt.url, c.matt.token, mattermost.Integrations{AddCommentFromApi: c.matt.integrations.AddCommentFromApi, DisableCalendarTaskNotificationApi: c.matt.integrations.DisableCalendarTaskNotificationApi, AllowedHosts: c.matt.integrations.AllowedHosts})
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

func (c *Container) GetOAuth2Authentik() *oauth2authentik.OAuth2 {
	return oauth2authentik.New(c.oAuth2Authentik.oAuth2Config, c.oAuth2Authentik.oidcProvider, c.oAuth2Authentik.logoutUrl)
}

func (c *Container) getADRepository() *ad.Repository {
	return ad.New(c.ads, c.adconfigs)
}
