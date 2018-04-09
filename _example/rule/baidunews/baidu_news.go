package baidunews

import (
	"github.com/nange/gospider/spider"
	"github.com/sirupsen/logrus"
)

func init() {
	spider.Register(rule)
}

var rule = &spider.TaskRule{
	Name:            "百度新闻规则",
	Description:     "抓取百度新闻各个分类的最新焦点新闻",
	Namespace:       "baidu_news",
	OutputFields:    []string{"category", "title", "link"},
	AllowURLRevisit: true,
	Rule: &spider.Rule{
		Head: func(ctx *spider.Context) error {
			return ctx.VisitForNext("http://news.baidu.com")
		},
		Nodes: map[int]*spider.Node{
			0: &spider.Node{ // 第一步: 获取所有分类
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Println("Visting", req.URL.String())
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement){
					`.menu-list a`: func(ctx *spider.Context, el *spider.HTMLElement) { // 获取所有分类
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
			1: &spider.Node{ // 第二步: 获取每个分类的新闻标题链接
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Println("Visting", req.URL.String())
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement){
					`#pane-news a`: func(ctx *spider.Context, el *spider.HTMLElement) {
						title := el.Text
						link := el.Attr("href")
						if title == "" || link == "javascript:void(0);" {
							return
						}
						category := el.Request.GetReqContextValue("category")
						ctx.Output(map[int]interface{}{
							0: category,
							1: title,
							2: link,
						})
					},
					`#col_focus a`: func(ctx *spider.Context, el *spider.HTMLElement) {
						title := el.Text
						link := el.Attr("href")
						if title == "" || link == "javascript:void(0);" {
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
