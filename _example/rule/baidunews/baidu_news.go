package baidunews

import (
	"github.com/nange/gospider/spider"
	log "github.com/sirupsen/logrus"
)

func init() {
	spider.Register(rule)
}

var (
	outputFields  = []string{"category", "title", "link"}
	outputFields2 = []string{"category", "category_link"}

	namespace1 = "baidu_news"
	namespace2 = "baidu_category"
)
var multNamespaceConf = map[string]*spider.MultipleNamespaceConf{
	namespace1: {
		OutputFields:      outputFields,
		OutputConstraints: spider.NewStringsConstraints(outputFields, 64, 128, 512),
	},
	namespace2: {
		OutputFields:      outputFields2,
		OutputConstraints: spider.NewStringsConstraints(outputFields2, 64, 256),
	},
}

// 演示如何在一条规则里面，同时需要导出数据到两张表
var rule = &spider.TaskRule{
	Name:                      "百度新闻规则",
	Description:               "抓取百度新闻各个分类的最新焦点新闻以及最新的新闻分类和链接",
	OutputToMultipleNamespace: true,
	MultipleNamespaceConf:     multNamespaceConf,
	Rule: &spider.Rule{
		Head: func(ctx *spider.Context) error {
			return ctx.VisitForNext("http://news.baidu.com")
		},
		Nodes: map[int]*spider.Node{
			0: step1, // 第一步: 获取所有分类
			1: step2, // 第二步: 获取每个分类的新闻标题链接
		},
	},
}

var step1 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Println("Visting", req.URL.String())
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`#channel-all .menu-list a`: func(ctx *spider.Context, el *spider.HTMLElement) error { // 获取所有分类
			category := el.Text
			ctx.PutReqContextValue("category", category)
			link := el.Attr("href")

			if category != "首页" {
				err := ctx.Output(map[int]interface{}{
					0: category,
					1: ctx.AbsoluteURL(link),
				}, namespace2)
				if err != nil {
					return err
				}
			}

			return ctx.VisitForNextWithContext(link)
		},
	},
}

var step2 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Println("Visting", req.URL.String())
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`#col_focus a, .focal-news a, .auto-col-focus a, .l-common .fn-c a`: func(ctx *spider.Context, el *spider.HTMLElement) error {
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
			}, namespace1)
		},
	},
}
