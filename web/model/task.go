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
	ID           uint64            `json:"id,string" ddb:"id"`
	TaskName     string            `json:"task_name" ddb:"task_name"`
	TaskRuleName string            `json:"task_rule_name" ddb:"task_rule_name"`
	TaskDesc     string            `json:"task_desc" ddb:"task_desc"`
	Status       common.TaskStatus `json:"status" ddb:"status"`
	Counts       int               `json:"counts" ddb:"counts"`
	// 参数配置部分
	OptUserAgent              string `json:"opt_user_agent" ddb:"opt_user_agent"`
	OptMaxDepth               int    `json:"opt_max_depth" ddb:"opt_max_depth"`
	OptAllowedDomains         string `json:"opt_allowed_domains" ddb:"opt_allowed_domains"`
	OptURLFilters             string `json:"opt_url_filters" ddb:"opt_url_filters"`
	OptAllowURLRevisit        bool   `json:"opt_allow_url_revisit" ddb:"opt_allow_url_revisit"`
	OptMaxBodySize            int    `json:"opt_max_body_size" ddb:"opt_max_body_size"`
	OptIgnoreRobotsTxt        bool   `json:"opt_ignore_robots_txt" ddb:"opt_ignore_robots_txt"`
	OptParseHTTPErrorResponse bool   `json:"opt_parse_http_error_response" ddb:"opt_parse_http_error_response"`
	// 频率限制
	LimitEnable       bool      `json:"limit_enable" ddb:"limit_enable"`
	LimitDomainRegexp string    `json:"limit_domain_regexp" ddb:"limit_domain_regexp"`
	LimitDomainGlob   string    `json:"limit_domain_glob" ddb:"limit_domain_glob"`
	LimitDelay        int       `json:"limit_delay" ddb:"limit_delay"`
	LimitRandomDelay  int       `json:"limit_random_delay" ddb:"limit_random_delay"`
	LimitParallelism  int       `json:"limit_parallelism" ddb:"limit_parallelism"`
	CreatedAt         time.Time `json:"created_at" ddb:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" ddb:"updated_at"`
}

func (o *Task) TableName() string {
	return "task"
}

func GetTaskList(db *gorm.DB, size, offset int) ([]Task, int, error) {
	queryset := NewTaskQuerySet(db)
	ret := make([]Task, 0)
	if err := queryset.All(&ret); err != nil {
		return nil, 0, errors.WithStack(err)
	}

	count, err := queryset.Count()
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return ret, count, nil
}
