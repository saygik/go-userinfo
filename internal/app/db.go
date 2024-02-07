package app

import (
	"database/sql"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/saygik/go-userinfo/internal/config"
)

func (a *App) newMsSQLConnect(cfg config.DBConfig) (*gorp.DbMap, error) {
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", cfg.Server, cfg.Dbname, cfg.User, cfg.Password))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqlServerDialect{}}
	return dbmap, nil
}

func (a *App) newGLPISQLConnect(cfg config.DBConfig) (*gorp.DbMap, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.User, cfg.Password, cfg.Server, cfg.Dbname))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqlServerDialect{}}
	return dbmap, nil
}
