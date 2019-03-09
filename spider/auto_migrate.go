package spider

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jinzhu/gorm"
)

// all code copied from gorm, just do some hack to support model defined by []string and map[string]constraints

// OutputConstraint is the output constraint of db
type OutputConstraint struct {
	SQL         string
	Index       string
	UniqueIndex string
}

// NewSQLString is the convenience func to return varchar sql string
func NewSQLString(size int, defaultValue ...string) (sql string) {
	if len(defaultValue) == 0 {
		sql = fmt.Sprintf("VARCHAR(%d) NOT NULL DEFAULT ''", size)
	} else {
		sql = fmt.Sprintf("VARCHAR(%d) NOT NULL DEFAULT '%s'", size, defaultValue[0])
	}
	return
}

// NewStringsConstraints is the convenience func to return varchar sql string of a batch columns
func NewStringsConstraints(columns []string, size ...int) (constraints map[string]*OutputConstraint) {
	s := make([]interface{}, len(size))
	for i, v := range size {
		s[i] = v
	}
	return NewConstraints(columns, s...)
}

// NewConstraints is the convenience func to return the custom constraints
func NewConstraints(columns []string, sizeOrSQLConstraint ...interface{}) (constraints map[string]*OutputConstraint) {
	constraints = make(map[string]*OutputConstraint)

	if len(columns) == 0 {
		panic("columns should contain at least 1 element")
	}

	switch len(sizeOrSQLConstraint) {
	case 0:
		panic("invalid parameter sizeOrSqlConstraint")
	case 1:
		switch v := sizeOrSQLConstraint[0].(type) {
		case int:
			sql := fmt.Sprintf("VARCHAR(%d) NOT NULL DEFAULT ''", v)

			for _, col := range columns {
				constraints[col] = &OutputConstraint{SQL: sql}
			}
		case string:
			if len(columns) > 1 {
				panic("default sizeOrSqlConstraint for all columns should be integer")
			} else {
				constraints[columns[0]] = &OutputConstraint{SQL: v}
			}
		default:
			panic("invalid parameter type")
		}
	default:
		if len(columns) != len(sizeOrSQLConstraint) {
			panic("length of column and sizeOrSqlConstraint are not match")
		}

		for idx, col := range columns {
			switch v := sizeOrSQLConstraint[idx].(type) {
			case int:
				constraints[col] = &OutputConstraint{SQL: fmt.Sprintf("VARCHAR(%d) NOT NULL DEFAULT ''", v)}
			case string:
				constraints[col] = &OutputConstraint{SQL: v}
			default:
				panic(fmt.Sprintf("parameter form idx<%d>, column<%s> is invalid", idx, col))
			}
		}
	}

	return
}

// AutoMigrateHack auto create table of the rule
func AutoMigrateHack(s *gorm.DB, rule *TaskRule) *gorm.DB {
	scope := s.NewScope(nil)
	s = autoMigrate(scope, rule).DB()

	return s
}

func autoMigrate(scope *gorm.Scope, rule *TaskRule) (s *gorm.Scope) {
	if rule.OutputToMultipleNamespace {
		keys := make([]string, 0, len(rule.MultipleNamespaceConf))
		for key := range rule.MultipleNamespaceConf {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			s = autoMigrateSingle(scope, rule, key)
		}
	} else {
		s = autoMigrateSingle(scope, rule, rule.Namespace)
	}

	return
}

func autoMigrateSingle(scope *gorm.Scope, rule *TaskRule, table string) *gorm.Scope {
	outputFields, constraints, _ := getOutputOpts(rule, table)
	if len(outputFields) < 1 {
		return scope
	}

	scope.Search.Table(table)
	tableName := scope.TableName()
	quotedTableName := scope.QuotedTableName()

	if !scope.Dialect().HasTable(tableName) {
		createTable(scope, rule, table)
	} else {
		for _, field := range outputFields {
			if !scope.Dialect().HasColumn(tableName, field) {
				sqlTag := getColumnTag(field, constraints)
				scope.Raw(fmt.Sprintf("ALTER TABLE %v ADD %v %v;", quotedTableName, scope.Quote(field), sqlTag)).Exec()
			}
		}
		autoIndex(scope, rule, table)
	}
	return scope
}

func createTable(scope *gorm.Scope, rule *TaskRule, table string) *gorm.Scope {
	var tags []string
	var primaryKeys []string

	foundID := false
	foundCreatedAt := false
	columns, constraints, outputTableOpts := getOutputOpts(rule, table)

	for _, field := range columns {
		isPrimaryKey := false
		lowerField := strings.ToLower(field)
		sqlTag := getColumnTag(field, constraints)

		if lowerField == "id" {
			foundID = true
		}

		if lowerField == "created_at" {
			foundCreatedAt = true
		}

		lowerSQLTag := strings.ToLower(sqlTag)
		if strings.Contains(lowerSQLTag, "primary key") {
			isPrimaryKey = true
			tags = append(tags, scope.Quote(field)+" "+strings.Replace(lowerSQLTag, "primary key", "", 1))
		} else {
			tags = append(tags, scope.Quote(field)+" "+sqlTag)
		}

		if isPrimaryKey {
			primaryKeys = append(primaryKeys, scope.Quote(field))
		}
	}

	if !foundID {
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

	scope.Raw(fmt.Sprintf("CREATE TABLE %v (%v %v)%s", scope.QuotedTableName(), strings.Join(tags, ","), primaryKeyStr, getTableOptions(outputTableOpts))).Exec()

	autoIndex(scope, rule, table)
	return scope
}

func getOutputOpts(rule *TaskRule, table string) (outputFields []string, outputConstraints map[string]*OutputConstraint, outputTableOpts string) {
	if rule.OutputToMultipleNamespace {
		outputFields = rule.MultipleNamespaceConf[table].OutputFields
		outputConstraints = rule.MultipleNamespaceConf[table].OutputConstraints
		outputTableOpts = rule.MultipleNamespaceConf[table].OutputTableOpts
	} else {
		outputFields = rule.OutputFields
		outputConstraints = rule.OutputConstraints
		outputTableOpts = rule.OutputTableOpts
	}

	return
}

func autoIndex(scope *gorm.Scope, rule *TaskRule, table string) *gorm.Scope {
	var indexes = map[string][]string{}
	var uniqueIndexes = map[string][]string{}

	cols, constraints, _ := getOutputOpts(rule, table)
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

func getTableOptions(outputTableOpts string) string {
	if outputTableOpts == "" {
		return " CHARSET=utf8mb4"
	}

	return " " + outputTableOpts
}

func getColumnTag(column string, constraints map[string]*OutputConstraint) (sqlTag string) {
	sqlTag = "varchar(255) NOT NULL DEFAULT ''"

	if constraints == nil {
		return
	}

	if c, ok := constraints[column]; ok {
		if c.SQL != "" {
			sqlTag = c.SQL
		}
	}

	return
}
