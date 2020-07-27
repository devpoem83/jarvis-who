package config

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"gitlab.eduwill.net/dev_team/jarvis-who/app/common"
	"strconv"
)

type DataSource struct {
	DbDriver                string
	DbServer                string
	DbName                  string
	DbUser                  string
	DbPassword              string
	DbOptionsPort           int
	DbOptionsReadOnlyIntent bool
	DbOptionsEncrypt        bool
}

var ganagosi_db *sql.DB
var log_db *sql.DB

func DbInit() {
	err1 := connectGanagosiDB()
	if err1 != nil {
		panic(err1)
	}

	err2 := connectLogDB()
	if err2 != nil {
		panic(err2)
	}
}

func connectGanagosiDB() error {
	var dataSource DataSource
	profile := GetProfile()

	dataSource.DbDriver = profile.GanagosiDbDriver
	dataSource.DbServer = profile.GanagosiDbServer
	dataSource.DbName = profile.GanagosiDbName
	dataSource.DbUser = profile.GanagosiDbUser
	dataSource.DbPassword = profile.GanagosiDbPassword
	dataSource.DbOptionsPort = profile.GanagosiDbOptionsPort
	dataSource.DbOptionsReadOnlyIntent = profile.GanagosiDbOptionsReadOnlyIntent
	dataSource.DbOptionsEncrypt = profile.GanagosiDbOptionsEncrypt

	return connectDb(&ganagosi_db, dataSource)
}

func connectLogDB() error {
	var dataSource DataSource
	profile := GetProfile()

	dataSource.DbDriver = profile.LogDbDriver
	dataSource.DbServer = profile.LogDbServer
	dataSource.DbName = profile.LogDbName
	dataSource.DbUser = profile.LogDbUser
	dataSource.DbPassword = profile.LogDbPassword
	dataSource.DbPassword = profile.LogDbPassword
	dataSource.DbOptionsPort = profile.LogDbOptionsPort
	dataSource.DbOptionsReadOnlyIntent = profile.LogDbOptionsReadOnlyIntent
	dataSource.DbOptionsEncrypt = profile.LogDbOptionsEncrypt

	return connectDb(&log_db, dataSource)
}

func GetGanagosiDB() *sql.DB {
	err := ganagosi_db.Ping()
	if err != nil {
		connectGanagosiDB()
	}
	return ganagosi_db
}

func GetLogDB() *sql.DB {
	err := log_db.Ping()
	if err != nil {
		connectGanagosiDB()
	}
	return log_db
}

func connectDb(db **sql.DB, datasource DataSource) error {
	dbDriver := "mssql"
	dbURL := "server=" + datasource.DbServer + ";user id=" + datasource.DbUser + ";password=" + datasource.DbPassword + ";database=" + datasource.DbName + ";port=" + strconv.Itoa(datasource.DbOptionsPort)
	common.Logger.Debug("dbURL : ", dbURL)
	common.Logger.Debug("DB Connection start... : ")

	var err error
	*db, err = sql.Open(dbDriver, dbURL)
	if err != nil || db == nil {
		common.Logger.Error("DB Connection - Fail!")
		panic(err.Error())
	} else {
		common.Logger.Debug("DB Check and create DB reference complete!")
	}

	return err
}
