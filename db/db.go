package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-gorp/gorp"
	"github.com/redis/go-redis/v9"
)

//_ "github.com/lib/pq" //import postgres

// DB ...
type DB struct {
	*sql.DB
}

var db *gorp.DbMap

// Init ...
func Init() {
	dbinfo := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", os.Getenv("DB_SERVER"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"))
	var err error
	db, err = ConnectDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}
}

// ConnectDB ...
func ConnectDB(dataSourceName string) (*gorp.DbMap, error) {
	db, err := sql.Open("mssql", dataSourceName)
	//	db, _ := sql.Open("mysql", "dellis:@/shud")
	//defer db.Close()
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected")
	if err = db.Ping(); err != nil {
		return nil, err
	}
	var version string
	db.QueryRow("SELECT @@VERSION").Scan(&version)
	fmt.Println("Connected to:", version)
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqlServerDialect{}}
	//dbmap.TraceOn("[gorp]", log.New(os.Stdout, "golang-gin:", log.Lmicroseconds)) //Trace database requests
	return dbmap, nil
}

// GetDB ...
func GetDB() *gorp.DbMap {
	return db
}

func CloseDB() {
	db.Db.Close()
}

var dbskype *gorp.DbMap

// InitDbSkype ...
func InitDbSkype() {
	dbinfo := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", os.Getenv("DB_SERVER_SKYPE"), os.Getenv("DB_NAME_SKYPE"), os.Getenv("DB_USER_SKYPE"), os.Getenv("DB_PASS_SKYPE"))
	var err error
	dbskype, err = ConnectDBSkype(dbinfo)
	if err != nil {
		log.Fatal(err)
	}
}

// ConnectDBSkype ...
func ConnectDBSkype(dataSourceName string) (*gorp.DbMap, error) {
	dbskype, err := sql.Open("mssql", dataSourceName)

	if err != nil {
		return nil, err
	}
	fmt.Println("Connected Db Skype")
	if err = dbskype.Ping(); err != nil {
		return nil, err
	}
	var version string
	dbskype.QueryRow("SELECT @@VERSION").Scan(&version)
	fmt.Println("Connected to:", version)
	dbmap := &gorp.DbMap{Db: dbskype, Dialect: gorp.SqlServerDialect{}}
	return dbmap, nil
}

// GetDBSkype ...
func GetDBSkype() *gorp.DbMap {
	return dbskype
}
func CloseDBSkype() {
	dbskype.Db.Close()
}

// RedisClient ...
var RedisClient *redis.Client

// InitRedis ...
func InitRedis(params ...string) {

	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

	db, _ := strconv.Atoi(params[0])

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       db,
	})
}

// GetRedis ...
func GetRedis() *redis.Client {
	return RedisClient
}

func PingRedis() error {
	var ctx = context.Background()
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}
