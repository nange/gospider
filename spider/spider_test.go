package spider

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/nange/gospider/common"
	"github.com/stretchr/testify/suite"
)

type testSpiderSuite struct {
	suite.Suite
	task Task
	ts   *httptest.Server
}

func (s *testSpiderSuite) SetupSuite() {
	task := Task{
		ID: 1,
		TaskRule: TaskRule{
			Name:                   "test_task",
			Namespace:              "test_table",
			OutputFields:           []string{"field1", "field2"},
			AllowURLRevisit:        true,
			IgnoreRobotsTxt:        true,
			ParseHTTPErrorResponse: true,
			InsecureSkipVerify:     true,
			DisableCookies:         true,
		},
		TaskConfig: TaskConfig{
			Option: Option{
				UserAgent: "gospider",
				MaxDepth:  10,
				//AllowedDomains: []string{"localhost", "127.0.0.1"},
				MaxBodySize:    100000,
				RequestTimeout: time.Second,
			},
			Limit: Limit{
				Enable:      true,
				DomainGlob:  "*",
				Delay:       time.Millisecond,
				RandomDelay: 100 * time.Millisecond,
				Parallelism: 2,
			},
			OutputConfig: OutputConfig{Type: common.OutputTypeMySQL},
		},
	}

	s.task = task
	rule := s.getTaskRule()
	s.task.Rule = rule

	s.ts = newTestServer()
}

func (s *testSpiderSuite) TearDownSuite() {
	s.ts.Close()
	if s.DirExists("./csv_output") {
		s.NoError(os.RemoveAll("./csv_output"))
	}
}

func (s *testSpiderSuite) TestRun() {
	task := s.task
	retChan := make(chan common.MTS, 1)

	db, mock, err := sqlmock.New()
	s.Require().NoError(err)
	defer db.Close()

	for i := 0; i < 4; i++ {
		mock.ExpectExec("(?i)insert into `test_table` (.+) values").
			WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	gs := New(&task, retChan)
	gs.SetDB(db)

	err = gs.Run()
	s.NoErrorf(err, "go spider run should no error")

	select {
	case ret := <-retChan:
		s.Equal(uint64(1), ret.ID)
		s.Equal(common.TaskStatusCompleted, ret.Status)
	case <-time.After(3 * time.Second):
		s.Fail("recive ret timeout")
	}

	err = mock.ExpectationsWereMet()
	s.NoError(err)

	task.OutputConfig.Type = common.OutputTypeCSV
	task.OutputConfig.CSVConf.CSVFilePath = "./csv_output"
	retChan2 := make(chan common.MTS, 1)
	gs = New(&task, retChan2)
	s.NoError(gs.Run())
	<-retChan2

}

func (s *testSpiderSuite) TestRunFail() {
	rule := &Rule{
		Head: func(ctx *Context) error {
			return errors.New("some error")
		},
		Nodes: map[int]*Node{
			0: {},
			1: {},
		},
	}
	task := s.task
	task.Rule = rule

	retChan := make(chan common.MTS, 1)
	gs := New(&task, retChan)
	s.NotNil(gs.Run())

	task = s.task
	task.TaskConfig.ProxyURLs = []string{"%!"}
	gs2 := New(&task, retChan)
	s.NotNil(gs2.Run())
}

func TestSpiderRun(t *testing.T) {
	suite.Run(t, new(testSpiderSuite))
}

func (s *testSpiderSuite) getTaskRule() *Rule {
	step1 := &Node{
		OnRequest: func(ctx *Context, req *Request) {
			s.T().Log("visiting:", req.URL)
		},
		OnResponse: func(ctx *Context, res *Response) error {
			s.T().Logf("code:%v", res.StatusCode)
			return nil
		},
		OnScraped: func(ctx *Context, res *Response) error {
			s.T().Log("scraped")
			return nil
		},
		OnHTML: map[string]func(ctx *Context, el *HTMLElement) error{
			`.category .item a`: func(ctx *Context, el *HTMLElement) error {
				category := el.Text
				link := el.Attr("href")
				ctx.PutReqContextValue("category", category)
				return ctx.VisitForNextWithContext(link)
			},
			`html`: func(ctx *Context, el *HTMLElement) error {
				childText := el.ChildText(`.item`)
				s.Equal("guoneiguoji", childText)
				childAttr := el.ChildAttr(`.item a`, "href")
				s.Equal("/guonei", childAttr)
				childAttrs := el.ChildAttrs(`.item a`, "href")
				s.ElementsMatch([]string{"/guonei", "/guoji"}, childAttrs)
				el.ForEach(`.item a`, func(i int, element *HTMLElement) {
					if i == 0 {
						s.Equal("guonei", element.Text)
					} else if i == 1 {
						s.Equal("guoji", element.Text)
					}
				})
				return nil
			},
		},
	}
	step2 := &Node{
		OnRequest: func(ctx *Context, req *Request) {
			s.T().Log("visiting:", req.URL)
		},
		OnResponse: func(ctx *Context, res *Response) error {
			s.T().Logf("code:%v", res.StatusCode)
			return nil
		},
		OnScraped: func(ctx *Context, res *Response) error {
			s.T().Log("scraped")
			return nil
		},
		OnHTML: map[string]func(ctx *Context, el *HTMLElement) error{
			`.news-item`: func(ctx *Context, el *HTMLElement) error {
				newsContent := el.Text
				category := ctx.GetReqContextValue("category")

				return ctx.Output(map[int]interface{}{
					0: category,
					1: newsContent,
				})
			},
		},
	}
	rule := &Rule{
		Head: func(ctx *Context) error {
			return ctx.VisitForNext(s.ts.URL)
		},
		Nodes: map[int]*Node{
			0: step1,
			1: step2,
		},
	}

	return rule
}

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>Test Page</title>
</head>
<body>
<h1>Hello World</h1>
<p class="description">This is a test page</p>
<p class="description">This is a test paragraph</p>
<div class="category">
<p class="item"><a href="/guonei">guonei</a></p>
<p class="item"><a href="/guoji">guoji</a></p>
</div>
</body>
</html>
		`))
	})
	mux.HandleFunc("/guonei", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>Test Page</title>
</head>
<body>
<p class="news-item">news item 1 of guonei</p>
<p class="news-item">news item 2 of guonei</p>
</body>
</html>
		`))
	})
	mux.HandleFunc("/guoji", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>Test Page</title>
</head>
<body>
<p class="news-item">news item 1 of guoji</p>
<p class="news-item">news item 2 of guoji</p>
</body>
</html>
		`))
	})

	return httptest.NewServer(mux)
}
