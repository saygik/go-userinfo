package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/saygik/go-userinfo/ad"
	"github.com/saygik/go-userinfo/controllers"
	"github.com/saygik/go-userinfo/db"
	ginlogrus "github.com/toorop/gin-logrus"
	"io/ioutil"

	//_ "github.com/saygik/go-glpi-api/log"
	"github.com/sirupsen/logrus"
	uuid "github.com/twinj/uuid"
	"net/http"
	"os"
	"runtime"
)

//CORSMiddleware ...
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

//RequestIDMiddleware ...
//Generate a unique ID and attach it to each request for future reference or use
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
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

	db.InitDbSkype()
	defer db.CloseDBSkype()
	//Start AD clients

	ad.Init(adconfig)
	defer ad.Close()
	//Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
	//Example: db.GetRedis().Set(KEY, VALUE, at.Sub(now)).Err()
	db.InitRedis("1")

	v1 := r.Group("/v1")
	{
		/*** START USER ***/
		user := new(controllers.UserController)

		controllers.DefaultDomain = os.Getenv("DEFAULT_DOMAIN")
		v1.GET("/users/ip", user.All)
		v1.GET("/users/ip/:domain", user.All)
		v1.GET("/user/ip", user.SetIp)
		v1.GET("/user/ip/:username", user.GetUserByName)

		aduser := new(controllers.ADUserController)
		v1.GET("/users/ad/:domain", aduser.All)
		v1.GET("/users/domains", aduser.AllDomains)

		skype := new(controllers.SkypeController)
		v1.GET("/skype/presences", skype.AllPresences)
		v1.GET("/skype/presence/:user", skype.OnePresence)
		v1.GET("/skype/activeconferences", skype.AllActiveConferences)
		v1.GET("/skype/conferencepresence/:id", skype.ConferencePresence)
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
