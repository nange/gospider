package core

import (
	"database/sql"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var modelList []Model

var db *gorm.DB

func SetGormDB(gdb *gorm.DB) {
	db = gdb
}

func GetGormDB() *gorm.DB {
	return db
}

func GetDB() *sql.DB {
	return db.DB()
}

type Model interface {
	TableName() string
}

func Register(model Model) {
	rv := reflect.ValueOf(model)
	if rv.IsNil() {
		panic("register model failed, model is nil")
	}
	for _, m := range modelList {
		if m.TableName() == model.TableName() {
			panic("register model failed, already have the table name:" + model.TableName())
		}
	}
	modelList = append(modelList, model)
}

func AutoMigrate() error {
	for _, model := range modelList {
		if err := db.Debug().Set("gorm:table_options", "CHARSET=utf8mb4").AutoMigrate(model).Error; err != nil {
			return errors.Wrap(err, "db auto migrate failed")
		}
	}

	return nil
}
