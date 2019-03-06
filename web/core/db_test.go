package core

import (
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testModel struct {
	ID   uint64 `gorm:"column:id;type:bigint unsigned AUTO_INCREMENT;primary_key"`
	Name string `gorm:"column:user_name;type:varchar(32) not null"`
}

func (t *testModel) TableName() string {
	return "gospider_testmodel"
}

func TestRegisterAndAutoMigrate(t *testing.T) {
	assert.Panics(t, func() {
		Register(nil)
	})
	assert.Panics(t, func() {
		var tm *testModel
		Register(tm)
	})

	Register(&testModel{})
	assert.Panics(t, func() {
		Register(&testModel{})
	})

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	gdb, err := gorm.Open("mysql", db)
	require.NoError(t, err)
	SetGormDB(gdb)

	mock.ExpectExec("(?i)create table `gospider_testmodel`").
		WillReturnResult(sqlmock.NewResult(0, 0))

	assert.NoError(t, AutoMigrate())
	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectExec("(?i)create table `gospider_testmodel`").
		WillReturnError(fmt.Errorf("some error"))
	assert.NotNil(t, AutoMigrate())
	assert.NoError(t, mock.ExpectationsWereMet())

}
