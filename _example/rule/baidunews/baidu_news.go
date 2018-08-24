package baidunews

import (
	"github.com/nange/gospider/spider"
	"github.com/sirupsen/logrus"
)

func init() {
	spider.Register(rule)
}

var outputFields = []string{"category", "title", "link"}

var rule = &spider.TaskRule{
	Name:         "百度新闻规则",
	Description:  "抓取百度新闻各个分类的最新焦点新闻",
	Namespace:    "baidu_news",
	OutputFields: outputFields,
	//OutputConstraints: map[string]*spider.OutputConstraint{
	//	outputFields[0]: &spider.OutputConstraint{Sql: "varchar(64) not null default ''"},
	//	outputFields[1]: &spider.OutputConstraint{Sql: "varchar(128) not null default ''"},
	//	outputFields[2]: &spider.OutputConstraint{Sql: "varchar(256) not null default ''"},
	//},
	OutputConstraints: spider.NewStringsConstraints(outputFields, 64, 128, 512), // 上面的简写方式
	AllowURLRevisit:   true,
	Rule: &spider.Rule{
		Head: func(ctx *spider.Context) error {
			return ctx.VisitForNext("http://news.baidu.com")
		},
		Nodes: map[int]*spider.Node{
			0: &spider.Node{ // 第一步: 获取所有分类
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Println("Visting", req.URL.String())
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
					`.menu-list a`: func(ctx *spider.Context, el *spider.HTMLElement) error { // 获取所有分类
						category := el.Text
						if category == "百家号" || category == "个性推荐" {
							return nil
						}
						if category == "首页" {
							category = "热点要闻"
						}

						ctx.PutReqContextValue("category", category)

						link := el.Attr("href")
						return ctx.VisitForNextWithContext(link)
					},
				},
			},
			1: &spider.Node{ // 第二步: 获取每个分类的新闻标题链接
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Println("Visting", req.URL.String())
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
					`#pane-news a`: func(ctx *spider.Context, el *spider.HTMLElement) error {
						title := el.Text
						link := el.Attr("href")
						if title == "" || link == "javascript:void(0);" {
							return nil
						}

						category := ctx.GetReqContextValue("category")
						return ctx.Output(map[int]interface{}{
							0: category,
							1: title,
							2: link,
						})
					},
					`#col_focus a`: func(ctx *spider.Context, el *spider.HTMLElement) error {
						title := el.Text
						link := el.Attr("href")
						if title == "" || link == "javascript:void(0);" {
							return nil
						}

						category := ctx.GetReqContextValue("category")
						return ctx.Output(map[int]interface{}{
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
