package spider

import (
	"crypto/tls"
	"net/http"
	"regexp"
	"time"

	// import the mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"github.com/nange/gospider/common"
	"github.com/pkg/errors"
)

// TaskConfig is the config of a task
type TaskConfig struct {
	CronSpec     string
	Option       Option
	Limit        Limit
	ProxyURLs    []string
	OutputConfig OutputConfig
}

// Option is the config option of a task
type Option struct {
	UserAgent              string
	MaxDepth               int
	AllowedDomains         []string
	URLFilters             []*regexp.Regexp
	AllowURLRevisit        bool
	MaxBodySize            int
	IgnoreRobotsTxt        bool
	InsecureSkipVerify     bool
	ParseHTTPErrorResponse bool
	DisableCookies         bool
	RequestTimeout         time.Duration
}

// Limit is the limit of a task
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

// OutputConfig is the output config of a task
type OutputConfig struct {
	Type      string
	CSVConf   CSVConf
	MySQLConf common.MySQLConf
}

// CSVConf is the csv conf of a task
type CSVConf struct {
	CSVFilePath string
}

func newCollector(config TaskConfig) (*colly.Collector, error) {
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

	if len(config.ProxyURLs) > 0 {
		rp, err := proxy.RoundRobinProxySwitcher(config.ProxyURLs...)
		if err != nil {
			return nil, errors.Wrapf(err, "set proxy switcher err")
		}
		c.SetProxyFunc(rp)
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

	if config.Option.RequestTimeout > 0 {
		c.SetRequestTimeout(config.Option.RequestTimeout)
	}

	if config.Option.InsecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.WithTransport(tr)
	}

	return c, nil
}
