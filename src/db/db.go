package db

import (
	"fmt"
	"time"
	"tinyquant/src/util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
)

var db *xorm.Engine

func GetSession() *xorm.Engine {

	return db
}

func InitMysql() {
	var err error

	db, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", util.MysqlUser, util.MysqlPass, util.MysqlHost, util.MysqlDBName))
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(60 * time.Second)
	//db.ShowSQL(true)
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

}

func CloseDB() {
	db.Close()
}
