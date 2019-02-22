package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nange/gospider/common"
	"github.com/nange/gospider/web/core"
	"github.com/pkg/errors"
)

//go:generate goqueryset -in task.go
// gen:qs
type Task struct {
	ID               uint64            `json:"id,string" gorm:"column:id;type:bigint unsigned AUTO_INCREMENT;primary_key"`
	TaskName         string            `json:"task_name" gorm:"column:task_name;type:varchar(64);not null;unique_index:uk_task_name"`
	TaskRuleName     string            `json:"task_rule_name" gorm:"column:task_rule_name;type:varchar(64);not null"`
	TaskDesc         string            `json:"task_desc" gorm:"column:task_desc;type:varchar(512);not null;default:''"`
	Status           common.TaskStatus `json:"status" gorm:"column:status;type:tinyint;not null;default:'0'"`
	Counts           int               `json:"counts" gorm:"column:counts;type:int;not null;default:'0'"`
	CronSpec         string            `json:"cron_spec" gorm:"column:cron_spec;type:varchar(64);not null;default:''"`
	OutputType       string            `json:"output_type" gorm:"column:output_type;type:varchar(64);not null;"`
	OutputExportDBID uint64            `json:"output_exportdb_id" gorm:"column:output_exportdb_id;type:bigint;not null;default:'0'"`
	// 参数配置部分
	OptUserAgent      string `json:"opt_user_agent" gorm:"column:opt_user_agent;type:varchar(128);not null;default:''"`
	OptMaxDepth       int    `json:"opt_max_depth" gorm:"column:opt_max_depth;type:int;not null;default:'0'"`
	OptAllowedDomains string `json:"opt_allowed_domains" gorm:"column:opt_allowed_domains;type:varchar(512);not null;default:''"`
	OptURLFilters     string `json:"opt_url_filters" gorm:"column:opt_url_filters;type:varchar(512);not null;default:''"`
	OptMaxBodySize    int    `json:"opt_max_body_size" gorm:"column:opt_max_body_size;type:int;not null;default:'0'"`
	OptRequestTimeout int    `json:"opt_request_timeout" gorm:"column:opt_request_timeout;type:int;not null;default:'10'"`
	// auto migrate
	AutoMigrate bool `json:"auto_migrate" gorm:"column:auto_migrate;type:tinyint;not null;default:'0'"`

	// 频率限制
	LimitEnable       bool      `json:"limit_enable" gorm:"column:limit_enable;type:tinyint;not null;default:'0'"`
	LimitDomainRegexp string    `json:"limit_domain_regexp" gorm:"column:limit_domain_regexp;type:varchar(128);not null;default:''"`
	LimitDomainGlob   string    `json:"limit_domain_glob" gorm:"column:limit_domain_glob;type:varchar(128);not null;default:''"`
	LimitDelay        int       `json:"limit_delay" gorm:"column:limit_delay;type:int;not null;default:'0'"`
	LimitRandomDelay  int       `json:"limit_random_delay" gorm:"column:limit_random_delay;type:int;not null;default:'0'"`
	LimitParallelism  int       `json:"limit_parallelism" gorm:"column:limit_parallelism;type:int;not null;default:'0'"`
	ProxyURLs         string    `json:"proxy_urls" gorm:"column:proxy_urls;type:varchar(2048);not null;default:''"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;index:idx_created_at"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;index:idx_updated_at"`
}

func (o *Task) TableName() string {
	return "gospider_task"
}

func init() {
	core.Register(&Task{})
}

func GetTaskList(db *gorm.DB, size, offset int) ([]Task, int, error) {
	queryset := NewTaskQuerySet(db)
	count, err := queryset.Count()
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	queryset = NewTaskQuerySet(db.Limit(size).Offset(offset))
	ret := make([]Task, 0)
	if err := queryset.OrderDescByCreatedAt().All(&ret); err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return ret, count, nil
}
