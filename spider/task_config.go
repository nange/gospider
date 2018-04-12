package spider

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
	"github.com/pkg/errors"
)

const (
	OutputTypeMySQL = "mysql"
	OutputTypeCSV   = "csv"
)

//TODO: 定时任务配置
type TaskConfig struct {
	CronSpec     string
	Option       Option
	Limit        Limit
	OutputConfig OutputConfig
}

type Option struct {
	UserAgent              string
	MaxDepth               int
	AllowedDomains         []string
	URLFilters             []*regexp.Regexp
	AllowURLRevisit        bool
	MaxBodySize            int
	IgnoreRobotsTxt        bool
	ParseHTTPErrorResponse bool
	DisableCookies         bool
}

type Limit struct {
	Enable bool
	// DomainRegexp is a regular expression to match against domains
	DomainRegexp string
	// DomainRegexp is a glob pattern to match against domains
	DomainGlob string
	// Delay is the duration to wait before creating a new request to the matching domains
	Delay time.Duration
	// RandomDelay is the extra randomized duration to wait added to Delay before creating a new request
	RandomDelay time.Duration
	// Parallelism is the number of the maximum allowed concurrent requests of the matching domains
	Parallelism int
}

type OutputConfig struct {
	Type      string
	CSVConf   CSVConf
	MySQLConf MySQLConf
}

type CSVConf struct {
	CSVFilePath string
	CSVFileName string
}

type MySQLConf struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func newCollector(config TaskConfig) *colly.Collector {
	opts := make([]func(*colly.Collector), 0)

	opts = append(opts, colly.Async(true))
	if config.Option.MaxDepth > 1 {
		opts = append(opts, colly.MaxDepth(config.Option.MaxDepth))
	}
	if len(config.Option.AllowedDomains) > 0 {
		opts = append(opts, colly.AllowedDomains(config.Option.AllowedDomains...))
	}
	if config.Option.AllowURLRevisit {
		opts = append(opts, colly.AllowURLRevisit())
	}
	if config.Option.IgnoreRobotsTxt {
		opts = append(opts, colly.IgnoreRobotsTxt())
	}
	if config.Option.MaxBodySize > 0 {
		opts = append(opts, colly.MaxBodySize(config.Option.MaxBodySize))
	}
	if config.Option.UserAgent != "" {
		opts = append(opts, colly.UserAgent(config.Option.UserAgent))
	}
	if config.Option.ParseHTTPErrorResponse {
		opts = append(opts, colly.ParseHTTPErrorResponse())
	}
	if len(config.Option.URLFilters) > 0 {
		opts = append(opts, colly.URLFilters(config.Option.URLFilters...))
	}

	c := colly.NewCollector(opts...)
	if config.Option.DisableCookies {
		c.DisableCookies()
	}

	if config.Limit.Enable {
		var limit colly.LimitRule
		if config.Limit.Delay > 0 {
			limit.Delay = config.Limit.Delay
		}
		if config.Limit.DomainGlob != "" {
			limit.DomainGlob = config.Limit.DomainGlob
		} else {
			limit.DomainGlob = "*"
		}
		if config.Limit.DomainRegexp != "" {
			limit.DomainRegexp = config.Limit.DomainRegexp
		}
		if config.Limit.Parallelism > 0 {
			limit.Parallelism = config.Limit.Parallelism
		}
		if config.Limit.RandomDelay > 0 {
			limit.RandomDelay = config.Limit.RandomDelay
		}

		c.Limit(&limit)
	}

	return c
}

func newDB(conf MySQLConf) (*sql.DB, error) {
	db, err := manager.New(conf.DBName, conf.User, conf.Password, conf.Host).
		Port(conf.Port).
		Set(
			manager.SetCharset("utf8"),
			manager.SetParseTime(true),
			manager.SetLoc("Local"),
		).Open(true)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	return db, nil
}
