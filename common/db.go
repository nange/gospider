package common

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// MySQLConf is the mysql conf
type MySQLConf struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

// NewGormDB return a new gorm db instance
func NewGormDB(conf MySQLConf) (*gorm.DB, error) {
	dsn := getDSNWithDB(conf)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		code, ok := GetSQLErrCode(err)
		if !ok {
			return nil, errors.WithStack(err)
		}
		if code == 1049 { // Database not exists
			if err := createDatabase(conf); err != nil {
				return nil, err
			}
		}
		db, err = gorm.Open("mysql", dsn)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if conf.MaxIdleConns == 0 {
		db.DB().SetMaxIdleConns(3)
	}
	if conf.MaxOpenConns == 0 {
		db.DB().SetMaxOpenConns(5)
	}
	if conf.MaxLifetime == 0 {
		db.DB().SetConnMaxLifetime(time.Hour)
	}

	return db, nil
}

// NewDB returns a new sql.DB instance
func NewDB(conf MySQLConf) (*sql.DB, error) {
	dsn := getDSNWithDB(conf)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)

	return db, nil
}

func getbaseDSN(conf MySQLConf) string {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port)

	return dsn
}

func getDSNWithDB(conf MySQLConf) string {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	return dsn
}

// GetSQLErrCode returns error code if err is a mysql error
func GetSQLErrCode(err error) (int, bool) {
	mysqlErr, ok := errors.Cause(err).(*mysql.MySQLError)
	if !ok {
		return -1, false
	}

	return int(mysqlErr.Number), true
}

func createDatabase(conf MySQLConf) error {
	dsn := getbaseDSN(conf)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return errors.WithStack(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + conf.DBName)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
