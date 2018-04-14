package model

import (
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var modelList []Model

var db *gorm.DB

func SetDB(gdb *gorm.DB) {
	db = gdb
}

func GetDB() *gorm.DB {
	return db
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
		if err := db.Debug().AutoMigrate(model).Error; err != nil {
			return errors.Wrap(err, "db auto migrate failed")
		}
	}

	return nil
}
