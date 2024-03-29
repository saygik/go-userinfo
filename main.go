package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/saygik/go-userinfo/ad"
	"github.com/saygik/go-userinfo/controllers"
	"github.com/saygik/go-userinfo/db"
	"github.com/saygik/go-userinfo/glpidb"
	"github.com/saygik/go-userinfo/sp"
	ginlogrus "github.com/toorop/gin-logrus"

	"net/http"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

// CORSMiddleware ...
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// RequestIDMiddleware ...
// Generate a unique ID and attach it to each request for future reference or use
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		suuid, _ := uuid.NewUUID() //uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", suuid.String())
		c.Next()
	}
}

var auth = new(controllers.AuthController)

// TokenAuthMiddleware ...
// JWT Authentication middleware attached to each request that needs to be authenitcated to validate the access_token in the header
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.TokenValid(c)
		c.Next()
	}
}

// TokenAuthMiddleware ...
// JWT Authentication middleware attached to each request that needs to be authenitcated to validate the access_token in the header
func UserFromTokenTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.UserFromToken(c)
		c.Next()
	}
}

func LoadConfiguration(file string) (ad.Config, error) {
	cfg := ad.Config{}
	configFile, err := ioutil.ReadFile(file)

	if err != nil {
		return cfg, err
	}
	json.Unmarshal(configFile, &cfg.ADS)
	//jsonParser := json.NewDecoder(configFile)
	//jsonParser.Decode(&config)
	return cfg, err
}

var log = logrus.New()

func main() {

	gin.SetMode(gin.ReleaseMode)

	src, err := os.OpenFile("api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("err", err)
	}
	log.Out = src
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.Info("------------------Starting programm-------------")
	log.Formatter = customFormatter
	r := gin.New()
	r.Use(ginlogrus.Logger(log), gin.Recovery())
	//	j := `{"123":"1.1.1.1","1234":"2.2.2.2"}`
	//Load the .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, please create one in the root directory")
	}

	adconfig, err := LoadConfiguration("adconfig.json")
	if err != nil || adconfig.ADS == nil {
		log.Fatal("Error loading adconfig.json file, please create one in the root directory")
	}

	r.Use(CORSMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	//Start MYSQL database
	//Example: db.GetDB() - More info in the models folder
	db.Init()
	defer db.CloseDB()

	glpidb.Init()
	defer glpidb.CloseDB()

	db.InitDbSkype()
	defer db.CloseDBSkype()
	//Start AD clients

	ad.Init(adconfig)
	defer ad.Close()
	//Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
	//Example: db.GetRedis().Set(KEY, VALUE, at.Sub(now)).Err()

	db.InitRedis("1")
	if err := db.PingRedis(); err != nil {
		log.Fatal("Redis not found")
	}
	//Start Sharepoint client
	sp.Init()

	aduser := new(controllers.ADUserController)
	aduser.GetAllDomainsUsers(true)
	// ticker := time.NewTicker(10 * time.Minute)
	// quit := make(chan struct{})
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			// do stuff
	// 			aduser.GetAllDomainsUsers(false)
	// 		case <-quit:
	// 			ticker.Stop()
	// 			return
	// 		}
	// 	}
	// }()
	v1 := r.Group("/v1")
	{

		user := new(controllers.UserController)
		v1.GET("/loginoauth", user.LoginOauth)
		v1.GET("/loginuser", TokenAuthMiddleware(), user.LoginUser)
		v1.GET("/token", user.LoginCallback)
		v1.POST("/login", user.Login)
		v1.GET("/logout", user.Logout)
		v1.POST("/token/refresh", auth.Refresh)
		userIP := new(controllers.UserIPController)

		controllers.DefaultDomain = os.Getenv("DEFAULT_DOMAIN")
		v1.GET("/users/ip", userIP.All)
		v1.GET("/users/ip/:domain", userIP.All)
		v1.GET("/user/ip", userIP.SetIp)
		v1.GET("/user/ip/:username", userIP.GetUserByName)
		v1.PUT("/user/avatar/:username", TokenAuthMiddleware(), userIP.UpdateUserAvatar)
		v1.GET("/user/activity/:username", TokenAuthMiddleware(), userIP.GetUserWeekActivity)
		v1.GET("/user/ad/:username", TokenAuthMiddleware(), aduser.GetUserByName)
		v1.GET("/cuser/whoami", TokenAuthMiddleware(), aduser.Find)
		v1.GET("/cuser/resources", TokenAuthMiddleware(), userIP.CurrentUserResources)

		v1.GET("/users/ad/:domain", TokenAuthMiddleware(), aduser.All)
		v1.GET("/users/allad", TokenAuthMiddleware(), aduser.AllAd)
		v1.GET("/users/short-ad-info", TokenAuthMiddleware(), aduser.AllAdUsersShort)
		v1.GET("/computers/allad", TokenAuthMiddleware(), aduser.AllAdComputers)

		v1.GET("/ad/stats/counts", aduser.AllAdCounts)
		v1.GET("/users/ad/:domain/:group", aduser.GroupUsers)
		v1.GET("/users/domains", TokenAuthMiddleware(), aduser.AllDomains)

		v1.GET("/schedules/:id", userIP.GetSchedule)
		v1.GET("/schedule/tasks/:idc", userIP.GetScheduleTasks)
		v1.POST("/schedule/task", userIP.AddScheduleTask)
		v1.DELETE("/schedule/task/:id", userIP.DelScheduleTask)
		v1.PUT("/schedule/task/:id", userIP.UpdateScheduleTask)

		skype := new(controllers.SkypeController)
		v1.GET("/skype/presences", skype.AllPresences)
		v1.GET("/skype/presence/:user", skype.OnePresence)
		v1.GET("/skype/activeconferences", skype.AllActiveConferences)
		v1.GET("/skype/conferencepresence/:id", skype.ConferencePresence)
		v1.GET("/skype/presences2", skype.AllPresences)

		sp_controller := new(controllers.SPController)
		v1.GET("/sp/zals", sp_controller.All)

		ideco_controller := new(controllers.IDECOController)
		v1.GET("/iutm/wlist", UserFromTokenTokenMiddleware(), ideco_controller.GetWiteList)
		mattermost_controller := new(controllers.MattermostController)
		v1.GET("/matt/users", mattermost_controller.GetAll)

		/**************** GLPI ************************/
		glpi_controller := new(controllers.GLPIController)
		v1.GET("/user/glpi/:username", UserFromTokenTokenMiddleware(), glpi_controller.GetUserByName)

		v1.GET("/softwares", TokenAuthMiddleware(), glpi_controller.GetSoftwares)
		v1.GET("/software/:id", TokenAuthMiddleware(), glpi_controller.GetSoftware)
		v1.GET("/software/:id/users", TokenAuthMiddleware(), glpi_controller.GetSoftwareUsers)
		v1.GET("/user/softwares/:username", glpi_controller.GetUserSoftwares)
		v1.DELETE("/user/software/:username/:id", glpi_controller.DelOneUserSoftware)
		v1.POST("/user/software/:username", glpi_controller.AddOneUserSoftware)
		v1.POST("/software/user/:software", glpi_controller.AddOneSoftwareUser)

		v1.GET("/glpi/sum/otkaz", glpi_controller.GetStatOtkazSum)                            // * Всего отказов //
		v1.GET("/glpi/otkazes", glpi_controller.GetOtkazes)                                   // * Все отказвы //
		v1.GET("/glpi/nctickets", TokenAuthMiddleware(), glpi_controller.GetTicketsNonClosed) // * Все незакрытые заявки //
		v1.GET("/glpi/ticket/:id", TokenAuthMiddleware(), glpi_controller.GetTicket)          // * Все незакрытые заявки //

		v1.GET("/users/glpi", TokenAuthMiddleware(), glpi_controller.GetUsers)
		v1.GET("/cuser/glpi", TokenAuthMiddleware(), glpi_controller.CurrentUserGLPI)

		v1.GET("/glpi/tickets/stats", glpi_controller.GetStatTickets)
		v1.GET("/glpi/tickets/statsdays", glpi_controller.GetStatTicketsDays)

		v1.GET("/glpi/failures/stats", glpi_controller.GetStatFailures)
		v1.GET("/glpi/regions/stats", glpi_controller.GetStatRegions)
		v1.GET("/glpi/stats/top10performers", glpi_controller.GetStatTop10Performers)
		v1.GET("/glpi/stats/top10iniciators", glpi_controller.GetStatTop10Iniciators)
		v1.GET("/glpi/stats/top10groups", glpi_controller.GetStatTop10Groups)
		v1.GET("/glpi/stats/periodcounts", glpi_controller.GetStatPeriodTicketsCounts)
		v1.GET("/glpi/stats/periodrequestypes", glpi_controller.GetStatPeriodRequestTypes)
		v1.GET("/glpi/stats/period-org-treemap", glpi_controller.GetStatPeriodOrgTreemap)
		v1.GET("/glpi/stats/period-regions-month-days", glpi_controller.GetStatPeriodRegionDayCounts)

		manuals_controller := new(controllers.ManualsController)
		v1.GET("/manuals/orgcodes", manuals_controller.AllOrgCodes)

	}

	r.LoadHTMLGlob("./public/html/*")

	r.Static("/public", "./public")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginBoilerplateVersion": "v1.03",
			"goVersion":             runtime.Version(),
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	fmt.Println("SSL", os.Getenv("SSL"))
	port := os.Getenv("PORT")

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Run(":" + port)
}
