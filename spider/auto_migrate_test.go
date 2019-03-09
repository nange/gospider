package spider

import (
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAutoMigrateHack(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	gdb, err := gorm.Open("mysql", db)
	require.NoError(t, err)

	mock.ExpectExec("(?i)create table `test_table` (.+)`id` (.+)`field1` (.+)`field2` (.+)`created_at` (.+)utf8mb4").
		WillReturnResult(sqlmock.NewResult(0, 0))

	rule := &TaskRule{
		Name:         "test_name",
		Namespace:    "test_table",
		OutputFields: []string{"field1", "field2"},
	}
	err = AutoMigrateHack(gdb, rule).Error
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectExec("(?i)create table `test_table`").
		WillReturnError(fmt.Errorf("some error"))
	err = AutoMigrateHack(gdb, rule).Error
	assert.NotNil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// multiple output table case
	rule.OutputToMultipleNamespace = true
	rule.MultipleNamespaceConf = map[string]*MultipleNamespaceConf{
		"test_mult_table": {
			OutputFields: []string{"mtable_field1", "mtable_field1"},
		},
		"test_mult_table2": {
			OutputFields: []string{"mtable_field1", "mtable_field1"},
		},
	}

	mock.ExpectExec("(?i)create table `test_mult_table`").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("(?i)create table `test_mult_table2`").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = AutoMigrateHack(gdb, rule).Error
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestNewSqlString(t *testing.T) {
	s1 := NewSQLString(10)
	assert.Equal(t, "VARCHAR(10) NOT NULL DEFAULT ''", s1)

	s2 := NewSQLString(10, "default string")
	assert.Equal(t, "VARCHAR(10) NOT NULL DEFAULT 'default string'", s2)
}

func TestNewStringsConstraints(t *testing.T) {
	cons := NewStringsConstraints([]string{"field1", "field2", "field3"}, 10, 20, 30)
	assert.Equal(t, 3, len(cons))
	assert.Equal(t, "VARCHAR(10) NOT NULL DEFAULT ''", cons["field1"].SQL)
	assert.Equal(t, "VARCHAR(20) NOT NULL DEFAULT ''", cons["field2"].SQL)
	assert.Equal(t, "VARCHAR(30) NOT NULL DEFAULT ''", cons["field3"].SQL)
}

func TestNewConstraints(t *testing.T) {
	cons := NewConstraints([]string{"field1", "field2"}, 10, 20)
	assert.Equal(t, 2, len(cons))
	assert.Equal(t, "VARCHAR(10) NOT NULL DEFAULT ''", cons["field1"].SQL)
	assert.Equal(t, "VARCHAR(20) NOT NULL DEFAULT ''", cons["field2"].SQL)

	cons2 := NewConstraints([]string{"field1", "field2"}, 10)
	assert.Equal(t, 2, len(cons2))
	assert.Equal(t, "VARCHAR(10) NOT NULL DEFAULT ''", cons2["field1"].SQL)
	assert.Equal(t, "VARCHAR(10) NOT NULL DEFAULT ''", cons2["field2"].SQL)

	cons3 := NewConstraints([]string{"field1"}, "VARCHAR(10) NOT NULL DEFAULT ''")
	assert.Equal(t, 1, len(cons3))
	assert.Equal(t, "VARCHAR(10) NOT NULL DEFAULT ''", cons3["field1"].SQL)

	assert.Panics(t, func() {
		NewConstraints([]string{"field1", "field2"}, "VARCHAR(10) NOT NULL DEFAULT ''")
	})
	assert.Panics(t, func() {
		NewConstraints([]string{"field1"}, map[string]string{})
	})
	assert.Panics(t, func() {
		NewConstraints([]string{"field1", "field2"})
	})
	assert.Panics(t, func() {
		NewConstraints([]string{"field1", "field2"}, 10, 20, 30)
	})

	cons4 := NewConstraints([]string{"field1", "field2", "field3"}, 10, "VARCHAR(20) NOT NULL DEFAULT ''", 30)
	assert.Equal(t, 3, len(cons4))
	assert.Equal(t, "VARCHAR(10) NOT NULL DEFAULT ''", cons4["field1"].SQL)
	assert.Equal(t, "VARCHAR(20) NOT NULL DEFAULT ''", cons4["field2"].SQL)
	assert.Equal(t, "VARCHAR(30) NOT NULL DEFAULT ''", cons4["field3"].SQL)

	assert.Panics(t, func() {
		NewConstraints([]string{"field1", "field2", "field3"}, 10, 20.0, 30)
	})

}
