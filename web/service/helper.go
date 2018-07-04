package service

import (
	"regexp"
	"strings"
	"time"

	"github.com/nange/gospider/spider"
	"github.com/nange/gospider/web/core"
	"github.com/nange/gospider/web/model"
	"github.com/pkg/errors"
)

func GetSpiderTaskByModel(task *model.Task) (*spider.Task, error) {
	rule, err := spider.GetTaskRule(task.TaskRuleName)
	if err != nil {
		return nil, err
	}

	var optAllowedDomains []string
	if task.OptAllowedDomains != "" {
		optAllowedDomains = strings.Split(task.OptAllowedDomains, ",")
	}
	var urlFiltersReg []*regexp.Regexp
	if task.OptURLFilters != "" {
		urlFilters := strings.Split(task.OptURLFilters, ",")
		for _, v := range urlFilters {
			reg, err := regexp.Compile(v)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			urlFiltersReg = append(urlFiltersReg, reg)
		}
	}

	sdb := model.SysDB{}
	query := model.NewSysDBQuerySet(core.GetDB())
	if err := query.IDEq(task.OutputSysDBID).One(&sdb); err != nil {
		return nil, errors.WithStack(err)
	}

	config := spider.TaskConfig{
		CronSpec: task.CronSpec,
		Option: spider.Option{
			UserAgent:              task.OptUserAgent,
			MaxDepth:               task.OptMaxDepth,
			AllowedDomains:         optAllowedDomains,
			URLFilters:             urlFiltersReg,
			AllowURLRevisit:        rule.AllowURLRevisit,
			MaxBodySize:            task.OptMaxBodySize,
			IgnoreRobotsTxt:        rule.IgnoreRobotsTxt,
			ParseHTTPErrorResponse: rule.ParseHTTPErrorResponse,
			DisableCookies:         rule.DisableCookies,
		},
		Limit: spider.Limit{
			Enable:      task.LimitEnable,
			DomainGlob:  task.LimitDomainGlob,
			Delay:       time.Duration(task.LimitDelay) * time.Millisecond,
			RandomDelay: time.Duration(task.LimitRandomDelay) * time.Millisecond,
			Parallelism: task.LimitParallelism,
		},
		ProxyURLs: strings.Split(task.ProxyURLs, ","),
		OutputConfig: spider.OutputConfig{
			Type: task.OutputType,
			MySQLConf: spider.MySQLConf{
				Host:     sdb.Host,
				Port:     sdb.Port,
				User:     sdb.User,
				Password: sdb.Password,
				DBName:   sdb.DBName,
			},
		},
	}

	return spider.NewTask(task.ID, *rule, config), nil
}
