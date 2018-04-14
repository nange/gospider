package common

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type MySQLConf struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	MaxIdleConns int
	MaxOpenConns int
}

func NewGormDB(conf MySQLConf) (*gorm.DB, error) {
	args := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
	db, err := gorm.Open("mysql", args)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if conf.MaxIdleConns == 0 {
		db.DB().SetMaxIdleConns(3)
	}
	if conf.MaxOpenConns == 0 {
		db.DB().SetMaxOpenConns(5)
	}

	return db, nil
}
