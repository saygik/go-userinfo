package db

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-gorp/gorp"
	_redis "github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
)

//_ "github.com/lib/pq" //import postgres

//DB ...
type DB struct {
	*sql.DB
}

var db *gorp.DbMap

//Init ...
func Init() {
	dbinfo := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", os.Getenv("DB_SERVER"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"))
	var err error
	db, err = ConnectDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}
}

//ConnectDB ...
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

//GetDB ...
func GetDB() *gorp.DbMap {
	return db
}
func CloseDB() {
	db.Db.Close()
}

//RedisClient ...
var RedisClient *_redis.Client

//InitRedis ...
func InitRedis(params ...string) {

	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

	db, _ := strconv.Atoi(params[0])

	RedisClient = _redis.NewClient(&_redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       db,
	})
}

//GetRedis ...
func GetRedis() *_redis.Client {
	return RedisClient
}
