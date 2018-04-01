package common

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type MySQLConf struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func newGormDB(conf MySQLConf) (*gorm.DB, error) {
	args := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
	db, err := gorm.Open("mysql", args)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	db.DB().SetMaxIdleConns(3)
	db.DB().SetMaxOpenConns(5)

	return db, nil
}

func GetGormDBFromEnv() (*gorm.DB, error) {
	host, ok := os.LookupEnv("GOSPIDER_DB_HOST")
	if !ok {
		host = "127.0.0.1"
	}
	var port int
	portstr, ok := os.LookupEnv("GOSPIDER_DB_PORT")
	if !ok {
		port = 3306
	} else {
		port64, err := strconv.ParseInt(portstr, 10, 64)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		port = int(port64)
	}
	user, ok := os.LookupEnv("GOSPIDER_DB_USER")
	if !ok {
		user = "root"
	}
	password, ok := os.LookupEnv("GOSPIDER_DB_PASSWORD")
	if !ok {
		password = ""
	}
	dbname, ok := os.LookupEnv("GOSPIDER_DB_NAME")
	if !ok {
		dbname = "test"
	}

	return newGormDB(MySQLConf{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbname,
	})
}
