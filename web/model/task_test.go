package model

import (
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
)

func (s *testModelSuite) TestGetTaskList() {
	rows1 := sqlmock.NewRows([]string{"count(*)"}).AddRow(10)
	s.mock.ExpectQuery("(?i)select count\\(\\*\\) from `gospider_task`").
		WillReturnRows(rows1)
	rows2 := sqlmock.NewRows([]string{"id"}).
		AddRow(uint64(1)).
		AddRow(uint64(2)).
		AddRow(uint64(3)).
		AddRow(uint64(4)).
		AddRow(uint64(5))
	s.mock.ExpectQuery("(?i)select \\* from `gospider_task`").
		WillReturnRows(rows2)

	list, count, err := GetTaskList(s.gdb, 5, 0)
	s.Require().NoError(err)
	s.Equal(10, count)
	s.Equal(5, len(list))
	s.Equal(uint64(1), list[0].ID)
	s.NoError(s.mock.ExpectationsWereMet())

	s.mock.ExpectQuery("(?i)select count\\(\\*\\) from `gospider_task`").
		WillReturnError(fmt.Errorf("some error"))
	list, count, err = GetTaskList(s.gdb, 5, 0)
	s.Require().NotNil(err)
	s.Equal(0, count)
	s.Nil(list)
	s.NoError(s.mock.ExpectationsWereMet())

	rows3 := sqlmock.NewRows([]string{"count(*)"}).AddRow(10)
	s.mock.ExpectQuery("(?i)select count\\(\\*\\) from `gospider_task`").
		WillReturnRows(rows3)
	s.mock.ExpectQuery("(?i)select \\* from `gospider_task` order by").
		WillReturnError(fmt.Errorf("some error"))
	list2, count2, err2 := GetTaskList(s.gdb, 5, 0)
	s.Require().NotNil(err2)
	s.Equal(0, count2)
	s.Nil(list2)
	s.NoError(s.mock.ExpectationsWereMet())
}
