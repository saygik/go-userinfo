package glpidb

import (
	"database/sql"
	"fmt"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

//_ "github.com/lib/pq" //import postgres

//DB ...
type DB struct {
	*sql.DB
}

var glpidb *gorp.DbMap

//Init ...
func Init() {
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", os.Getenv("GLPI_DB_USER"), os.Getenv("GLPI_DB_PASS"), os.Getenv("GLPI_DB_SERVER"), os.Getenv("GLPI_DB_NAME"))
	var err error
	glpidb, err = ConnectDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}
}

//ConnectDB ...
func ConnectDB(dataSourceName string) (*gorp.DbMap, error) {
	glpidb, err := sql.Open("mysql", dataSourceName)
	//	db, _ := sql.Open("mysql", "dellis:@/shud")
	//defer db.Close()
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected")
	if err = glpidb.Ping(); err != nil {
		return nil, err
	}
	var version string
	glpidb.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)
	dbmap := &gorp.DbMap{Db: glpidb, Dialect: gorp.MySQLDialect{}}
	//dbmap.TraceOn("[gorp]", log.New(os.Stdout, "golang-gin:", log.Lmicroseconds)) //Trace database requests
	return dbmap, nil
}

//GetDB ...
func GetDB() *gorp.DbMap {
	return glpidb
}
func CloseDB() {
	glpidb.Db.Close()
}
