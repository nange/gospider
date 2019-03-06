package model

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type testModelSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	gdb  *gorm.DB
}

func (s *testModelSuite) SetupSuite() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)
	s.db = db
	s.mock = mock
	gdb, err := gorm.Open("mysql", db)
	s.Require().NoError(err)
	s.gdb = gdb
}

func (s *testModelSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

func TestModel(t *testing.T) {
	suite.Run(t, new(testModelSuite))
}
