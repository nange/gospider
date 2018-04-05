package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nange/gospider/web/common"
	"github.com/pkg/errors"
)

//go:generate goqueryset -in task.go
// gen:qs
type Task struct {
	ID           uint64            `json:"id,string"`
	TaskName     string            `json:"task_name"`
	TaskRuleName string            `json:"task_rule_name"`
	TaskDesc     string            `json:"task_desc"`
	Status       common.TaskStatus `json:"status"`
	Counts       int               `json:"counts"`
	// 参数配置部分
	OptUserAgent              string `json:"opt_user_agent"`
	OptMaxDepth               int    `json:"opt_max_depth"`
	OptAllowedDomains         string `json:"opt_allowed_domains"`
	OptURLFilters             string `json:"opt_url_filters"`
	OptAllowURLRevisit        bool   `json:"opt_allow_url_revisit"`
	OptMaxBodySize            int    `json:"opt_max_body_size"`
	OptIgnoreRobotsTxt        bool   `json:"opt_ignore_robots_txt"`
	OptParseHTTPErrorResponse bool   `json:"opt_parse_http_error_response"`
	// 频率限制
	LimitEnable       bool      `json:"limit_enable"`
	LimitDomainRegexp string    `json:"limit_domain_regexp"`
	LimitDomainGlob   string    `json:"limit_domain_glob"`
	LimitDelay        int       `json:"limit_delay"`
	LimitRandomDelay  int       `json:"limit_random_delay"`
	LimitParallelism  int       `json:"limit_parallelism"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (o *Task) TableName() string {
	return "task"
}

func GetTaskList(db *gorm.DB, size, offset int) ([]Task, int, error) {
	queryset := NewTaskQuerySet(db)
	count, err := queryset.Count()
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	if size > 0 {
		db = db.Limit(size)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}
	queryset = NewTaskQuerySet(db)
	ret := make([]Task, 0)
	if err := queryset.OrderDescByCreatedAt().All(&ret); err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return ret, count, nil
}
