package spider

import (
	"log"
	"testing"
	"time"
)

func init() {
	Register(getRule())
}

func TestRun(t *testing.T) {
	config := getConfig()
	task := NewTask(*rules["百度热点要闻"], config)
	if err := Run(task); err != nil {
		t.Errorf("run task failed, err:%#v", err.Error())
		return
	}
	time.Sleep(10 * time.Second)
}

func getConfig() TaskConfig {
	return TaskConfig{
		Option: Option{
			MaxDepth:        1,
			AllowedDomains:  []string{"news.baidu.com"},
			IgnoreRobotsTxt: true,
			AllowURLRevisit: true,
		},
		Limit: Limit{
			DomainRegexp: "news.baidu.com",
			Parallelism:  1,
			//Delay:        time.Second,
		},
	}
}

func getRule() *TaskRule {
	name := "百度热点要闻"
	return &TaskRule{
		Name:         name,
		Namespace:    "baidu_news",
		Description:  "获取百度最新热点要闻标题和链接",
		OutputFields: []string{"category", "title", "link"},
		Rule: &Rule{
			Head: func(ctx *Context) error {
				return ctx.VisitForNext("http://news.baidu.com")
			},
			Nodes: map[int]*Node{
				0: &Node{
					OnRequest: func(ctx *Context, req *Request) {
						log.Println("Visting", req.URL.String())
					},
					OnError: func(ctx *Context, res *Response, err error) {
						log.Println("Visting err:", err.Error())
					},
					OnHTML: map[string]func(*Context, *HTMLElement){
						`.menu-list a`: func(ctx *Context, el *HTMLElement) { // 获取所有分类
							category := el.Text
							if category == "百家号" || category == "个性推荐" {
								return
							}
							if category == "首页" {
								category = "热点要闻"
							}

							el.Request.PutReqContextValue("category", category)

							link := el.Attr("href")
							el.Request.VisitForNextWithContext(link)
						},
					},
				},
				1: &Node{
					OnScraped: func(ctx *Context, res *Response) {
						log.Println("Scraped==================", res.Request.URL.String())
					},
					OnRequest: func(ctx *Context, req *Request) {
						log.Println("Visting", req.URL.String())
					},
					OnError: func(ctx *Context, res *Response, err error) {
						log.Println("Visting err:", err.Error())
					},
					OnHTML: map[string]func(*Context, *HTMLElement){
						`#pane-news a`: func(ctx *Context, el *HTMLElement) {
							title := el.Text
							link := el.Attr("href")
							if link == "javascript:void(0);" {
								return
							}
							category := el.Request.GetReqContextValue("category")
							ctx.Output(map[int]interface{}{
								0: category,
								1: title,
								2: link,
							})
						},
						`#col_focus a`: func(ctx *Context, el *HTMLElement) {
							title := el.Text
							if title == "" {
								return
							}
							link := el.Attr("href")
							if link == "javascript:void(0);" {
								return
							}
							category := el.Request.GetReqContextValue("category")
							ctx.Output(map[int]interface{}{
								0: category,
								1: title,
								2: link,
							})
						},
					},
				},
			},
		},
	}
}
