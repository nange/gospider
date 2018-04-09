package dianping

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nange/gospider/spider"
	"github.com/sirupsen/logrus"
)

func init() {
	spider.Register(rule)
}

// NOTICE: 目前这个例子仅实现了抓取美食类商家
var rule = &spider.TaskRule{
	Name:            "大众点评商家数据",
	Description:     "抓取大众点评上全国各大城市所有类型的商家详情数据",
	Namespace:       "dianping_shop",
	OutputFields:    []string{"city", "adname", "big_category", "sub_category", "shop_name", "address", "tel", "photos"},
	AllowURLRevisit: true,
	Rule: &spider.Rule{
		Head: func(ctx *spider.Context) error { // 定义入口
			return ctx.VisitForNext("http://www.dianping.com/citylist")
		},
		Nodes: map[int]*spider.Node{
			0: &spider.Node{ // 第一步: 获取所有城市的链接
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Infof("Visiting %s", req.URL.String())
				},
				OnError: func(ctx *spider.Context, res *spider.Response, err error) {
					logrus.Errorf("Visiting failed! url:%s, err:%s", res.Request.URL.String(), err.Error())
					// 出错时重试三次
					Retry(res.Request, 3)
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement){
					`.main-citylist .letter-item a`: func(ctx *spider.Context, el *spider.HTMLElement) {
						city := el.Text

						el.Request.PutReqContextValue("city", city)
						link := el.Attr("href")
						el.Request.VisitForNextWithContext(link)
					},
				},
			},
			1: &spider.Node{ // 第二步: 获取美食分类id
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Infof("Visiting %s", req.URL.String())
				},
				OnError: func(ctx *spider.Context, res *spider.Response, err error) {
					logrus.Errorf("Visiting failed! url:%s, err:%s", res.Request.URL.String(), err.Error())
					// 出错时重试三次
					Retry(res.Request, 3)
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement){
					`.first-cate .first-item .span-container`: func(ctx *spider.Context, el *spider.HTMLElement) {
						link := el.ChildAttr(".index-item", "href")
						if link == "" {
							return
						}
						trim := strings.TrimPrefix(link, el.Request.URL.String())
						items := strings.SplitN(trim, "/", 3)
						if len(items) < 2 {
							return
						}
						cateID := items[1]

						bigCate := el.ChildText(".index-title")
						if bigCate != "美食" {
							return
						}

						el.Request.PutReqContextValue("big_category", bigCate)
						el.Request.VisitForNextWithContext(el.Request.URL.String() + "/" + cateID)
					},
				},
			},
			2: &spider.Node{ // 第三步: 获取所有子分类
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Infof("Visiting %s", req.URL.String())
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement){
					`#classfy a`: func(ctx *spider.Context, el *spider.HTMLElement) {
						link := el.Attr("href")
						if link == "javascript:;" {
							return
						}
						subCate := el.Text
						el.Request.PutReqContextValue("sub_category", subCate)

						el.Request.VisitForNextWithContext(link)
					},
				},
			},
			3: &spider.Node{ // 第四步: 获取最小分类
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Infof("Visiting %s", req.URL.String())
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement){
					`.navigation`: func(ctx *spider.Context, el *spider.HTMLElement) {
						el.ForEach(`#classfy-sub a`, func(i int, element *spider.HTMLElement) {
							if i == 0 { // 第一个链接为"不限", 忽略
								return
							}
							el.Request.VisitForNextWithContext(element.Attr("href"))
						})

					},
				},
			},
			4: &spider.Node{ // 第五步: 使用最小分类和行政区域构造请求列表
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Infof("Visiting %s", req.URL.String())
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement){
					`#region-nav a`: func(ctx *spider.Context, el *spider.HTMLElement) {
						adname := el.Text
						el.Request.PutReqContextValue("adname", adname)

						el.Request.VisitForNextWithContext(el.Attr("href"))
					},
				},
			},
			5: &spider.Node{ // 第六步: 获取最大分页数并依次请求每一页数据
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Infof("Visiting %s", req.URL.String())
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement){
					`.page a:nth-last-child(2)`: func(ctx *spider.Context, el *spider.HTMLElement) {
						countTxt := el.Text
						if countTxt == "" {
							return
						}

						count64, err := strconv.ParseInt(countTxt, 10, 64)
						if err != nil {
							logrus.Errorf("pase page count err:%s", err.Error())
							return
						}

						for i := 2; i <= int(count64); i++ {
							nextURL := fmt.Sprintf("%sp%d", el.Request.URL.String(), i)
							logrus.Infof("nextURL:%s", nextURL)
							el.Request.Visit(nextURL)
						}
					},
					`.shop-list li .pic a`: func(ctx *spider.Context, el *spider.HTMLElement) {
						photo := el.ChildAttr(`img`, "src")
						el.Request.PutReqContextValue("photos", photo)

						el.Request.VisitForNextWithContext(el.Attr("href"))
					},
				},
			},
			6: &spider.Node{ // 第七步: 归总所有字段数据并导出
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Infof("Visiting %s", req.URL.String())
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement){
					`#basic-info`: func(ctx *spider.Context, el *spider.HTMLElement) {
						shopNameNodes := el.DOM.Find(`.shop-name`).Nodes
						if len(shopNameNodes) == 0 {
							return
						}
						shopName := shopNameNodes[0].FirstChild.Data
						address := el.ChildText(`.address .item`)
						tel := el.ChildText(`.tel .item`)

						city := el.Request.GetReqContextValue("city")
						adname := el.Request.GetReqContextValue("adname")
						bigCate := el.Request.GetReqContextValue("big_category")
						subCate := el.Request.GetReqContextValue("sub_category")
						photos := el.Request.GetReqContextValue("photos")

						ctx.Output(map[int]interface{}{
							0: city,
							1: adname,
							2: bigCate,
							3: subCate,
							4: shopName,
							5: address,
							6: tel,
							7: photos,
						})
					},
				},
			},
		},
	},
}

func Retry(req *spider.Request, count int) error {
	key := fmt.Sprintf("err_req_%s", req.URL.String())

	var et int
	if errCount := req.GetAnyReqContextValue(key); errCount != nil {
		et = errCount.(int)
		if et >= count {
			return fmt.Errorf("exceed %d counts", count)
		}
	}
	logrus.Infof("errCount:%d, we wil retry url:%s, after 1 second", et+1, req.URL.String())
	time.Sleep(time.Second)
	req.PutReqContextValue(key, et+1)
	req.Retry()

	return nil
}
