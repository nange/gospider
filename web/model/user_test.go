package model

import (
	"github.com/DATA-DOG/go-sqlmock"
)

func (s *testModelSuite) TestGenUserHashPassword() {
	hash, err := GenUserHashPassword("admin")
	s.NoErrorf(err, "should gen hash password success")
	s.T().Logf("hash:%s", hash)
}

func (s *testModelSuite) TestInitAdminUserIfNeeded() {
	row := sqlmock.NewRows([]string{"id"}).AddRow(uint64(1))
	s.mock.ExpectQuery("(?i)select \\* from `gospider_user` where \\(user_name = \\?\\)").
		WithArgs("admin").
		WillReturnRows(row)
	err := InitAdminUserIfNeeded(s.gdb)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())

	s.mock.ExpectQuery("(?i)select \\* from `gospider_user` where \\(user_name = \\?\\)").
		WithArgs("admin").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	s.mock.ExpectExec("(?i)insert into `gospider_user` (.+) values").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = InitAdminUserIfNeeded(s.gdb)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())

}

func (s *testModelSuite) TestIsValidUser() {
	pw, err := GenUserHashPassword("admin")
	s.Require().NoError(err)

	s.mock.ExpectQuery("(?i)select \\* from `gospider_user` where \\(user_name = \\?\\)").
		WithArgs("admin").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	valid, user, err := IsValidUser(s.gdb, "admin", pw)
	s.Nil(err)
	s.False(valid)
	s.Nil(user)
	s.NoError(s.mock.ExpectationsWereMet())

	row := sqlmock.NewRows([]string{"id", "password"}).AddRow(uint64(1), pw)
	s.mock.ExpectQuery("(?i)select \\* from `gospider_user` where \\(user_name = \\?\\)").
		WithArgs("admin").
		WillReturnRows(row)
	valid2, user2, err2 := IsValidUser(s.gdb, "admin", "admin")
	s.NoError(err2)
	s.True(valid2)
	s.NotNil(user2)
	s.NoError(s.mock.ExpectationsWereMet())

}
