package spider

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// all code copied from gorm, just do some hack to support model defined by []string and map[string]constraints

type OutputConstraint struct {
	Sql         string
	Index       string
	UniqueIndex string
}

func NewSqlString(size int, defaultValue ...string) (sql string) {
	if len(defaultValue) == 0 {
		sql = fmt.Sprintf("VARCHAR(%d) NOT NULL DEFAULT ''", size)
	} else {
		sql = fmt.Sprintf("VARCHAR(%d) NOT NULL DEFAULT '%s'", size, defaultValue[0])
	}
	return
}

func NewStringsConstraints(columns []string, size ...int) (constraints map[string]*OutputConstraint) {
	s := make([]interface{}, len(size))
	for i, v := range size {
		s[i] = v
	}
	return NewConstraints(columns, s...)
}

func NewConstraints(columns []string, strSizeOrSqlConstraint ...interface{}) (constraints map[string]*OutputConstraint) {
	constraints = make(map[string]*OutputConstraint)

	if len(columns) == 0 {
		logrus.Error("columns should contain at least 1 element")
		return
	}

	switch len(strSizeOrSqlConstraint) {
	case 0:
		logrus.Error("invalid parameter strSizeOrSqlConstraint")
		return
	case 1:
		switch v := strSizeOrSqlConstraint[0].(type) {
		case int:
			sql := fmt.Sprintf("VARCHAR(%d) NOT NULL DEFAULT ''", v)

			for _, col := range columns {
				constraints[col] = &OutputConstraint{Sql: sql}
			}
		case string:
			if len(columns) > 1 {
				logrus.Error("default strSizeOrSqlConstraint for all columns should be integer")
			} else {
				constraints[columns[0]] = &OutputConstraint{Sql: v}
			}
		default:
			logrus.Error("invalid parameter type")
		}
	default:
		if len(columns) != len(strSizeOrSqlConstraint) {
			logrus.Error("length of column and strSizeOrSqlConstraint are not match")
			return
		}

		for idx, col := range columns {
			switch v := strSizeOrSqlConstraint[idx].(type) {
			case int:
				constraints[col] = &OutputConstraint{Sql: fmt.Sprintf("VARCHAR(%d) NOT NULL DEFAULT ''", v)}
			case string:
				constraints[col] = &OutputConstraint{Sql: v}
			default:
				logrus.Error(fmt.Sprintf("parameter form idx<%d>, column<%s> is invalid", idx, col))
				return
			}
		}
	}

	return
}

func AutoMigrateHack(s *gorm.DB, rule *TaskRule) *gorm.DB {
	scope := s.NewScope(nil)
	s = autoMigrate(scope, rule).DB()

	return s
}

func autoMigrate(scope *gorm.Scope, rule *TaskRule) *gorm.Scope {
	columns := rule.OutputFields
	constraints := rule.OutputConstaints

	scope.Search.Table(rule.Namespace)
	tableName := scope.TableName()
	quotedTableName := scope.QuotedTableName()

	if !scope.Dialect().HasTable(tableName) {
		createTable(scope, rule)
	} else {
		for _, field := range columns {
			if !scope.Dialect().HasColumn(tableName, field) {
				sqlTag := getColumnTag(field, constraints)
				scope.Raw(fmt.Sprintf("ALTER TABLE %v ADD %v %v;", quotedTableName, scope.Quote(field), sqlTag)).Exec()
			}
		}
		autoIndex(scope, rule)
	}
	return scope
}

func createTable(scope *gorm.Scope, rule *TaskRule) *gorm.Scope {
	var tags []string
	var primaryKeys []string

	foundId := false
	foundCreatedAt := false

	columns := rule.OutputFields
	constraints := rule.OutputConstaints

	for _, field := range columns {
		isPrimaryKey := false
		lowerField := strings.ToLower(field)
		sqlTag := getColumnTag(field, constraints)

		if lowerField == "id" {
			foundId = true
		}

		if lowerField == "created_at" {
			foundCreatedAt = true
		}

		lowerSqlTag := strings.ToLower(sqlTag)
		if strings.Contains(lowerSqlTag, "primary key") {
			isPrimaryKey = true
			tags = append(tags, scope.Quote(field)+" "+strings.Replace(lowerSqlTag, "primary key", "", 1))
		} else {
			tags = append(tags, scope.Quote(field)+" "+sqlTag)
		}

		if isPrimaryKey {
			primaryKeys = append(primaryKeys, scope.Quote(field))
		}
	}

	if !foundId {
		tags = append([]string{"`id` bigint(64) unsigned NOT NULL AUTO_INCREMENT"}, tags...)
		primaryKeys = append(primaryKeys, `id`)
	}
	if !foundCreatedAt {
		tags = append(tags, "`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP")
	}

	var primaryKeyStr string
	if len(primaryKeys) > 0 {
		primaryKeyStr = fmt.Sprintf(", PRIMARY KEY (%v)", strings.Join(primaryKeys, ","))
	}

	scope.Raw(fmt.Sprintf("CREATE TABLE %v (%v %v)%s", scope.QuotedTableName(), strings.Join(tags, ","), primaryKeyStr, getTableOptions(rule))).Exec()

	autoIndex(scope, rule)
	return scope
}

func autoIndex(scope *gorm.Scope, rule *TaskRule) *gorm.Scope {
	var indexes = map[string][]string{}
	var uniqueIndexes = map[string][]string{}

	cols := rule.OutputFields
	constraints := rule.OutputConstaints

	if constraints == nil {
		return scope
	}

	for _, field := range cols {
		entry, ok := constraints[field]
		if !ok {
			continue
		}

		name := entry.Index
		if name != "" {
			names := strings.Split(name, ",")

			for _, name := range names {
				if name == "INDEX" || name == "" {
					name = scope.Dialect().BuildKeyName("idx", scope.TableName(), field)
				}
				indexes[name] = append(indexes[name], field)
			}
		}

		name = entry.UniqueIndex
		if name != "" {
			names := strings.Split(name, ",")

			for _, name := range names {
				if name == "UNIQUE_INDEX" || name == "" {
					name = scope.Dialect().BuildKeyName("uix", scope.TableName(), field)
				}
				uniqueIndexes[name] = append(uniqueIndexes[name], field)
			}
		}
	}

	for name, columns := range indexes {
		if db := scope.NewDB().Table(scope.TableName()).Model(scope.Value).AddIndex(name, columns...); db.Error != nil {
			scope.DB().AddError(db.Error)
		}
	}

	for name, columns := range uniqueIndexes {
		if db := scope.NewDB().Table(scope.TableName()).Model(scope.Value).AddUniqueIndex(name, columns...); db.Error != nil {
			scope.DB().AddError(db.Error)
		}
	}

	return scope
}

func getTableOptions(rule *TaskRule) string {
	if rule.OutputTableOpts == "" {
		return " CHARSET=utf8mb4"
	}

	return " " + rule.OutputTableOpts
}

func getColumnTag(column string, constraints map[string]*OutputConstraint) (sqlTag string) {
	sqlTag = "varchar(255) NOT NULL DEFAULT ''"

	if constraints == nil {
		return
	}

	if c, ok := constraints[column]; ok {
		if c.Sql != "" {
			sqlTag = c.Sql
		}
	}

	return
}