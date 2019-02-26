package dianping

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nange/gospider/spider"
	log "github.com/sirupsen/logrus"
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
			0: step1, // 第一步: 获取所有城市的链接
			1: step2, // 第二步: 获取美食分类id
			2: step3, // 第三步: 获取所有子分类
			3: step4, // 第四步: 获取最小分类,
			4: step5, // 第五步: 使用最小分类和行政区域构造请求列表
			5: step6, // 第六步: 获取最大分页数并依次请求每一页数据
			6: step7, // 第七步: 归总所有字段数据并导出
		},
	},
}

var step1 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Infof("Visiting %s", req.URL.String())
	},
	OnError: func(ctx *spider.Context, res *spider.Response, err error) error {
		log.Errorf("Visiting failed! url:%s, err:%s", res.Request.URL.String(), err.Error())
		// 出错时重试三次
		return Retry(ctx, 3)
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`.main-citylist .letter-item a`: func(ctx *spider.Context, el *spider.HTMLElement) error {
			city := el.Text

			ctx.PutReqContextValue("city", city)
			link := el.Attr("href")
			return ctx.VisitForNextWithContext(link)
		},
	},
}

var step2 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Infof("Visiting %s", req.URL.String())
	},
	OnError: func(ctx *spider.Context, res *spider.Response, err error) error {
		log.Errorf("Visiting failed! url:%s, err:%s", res.Request.URL.String(), err.Error())
		// 出错时重试三次
		return Retry(ctx, 3)
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`.first-cate .first-item .span-container`: func(ctx *spider.Context, el *spider.HTMLElement) error {
			link := el.ChildAttr(".index-item", "href")
			if link == "" {
				return nil
			}
			trim := strings.TrimPrefix(link, el.Request.URL.String())
			items := strings.SplitN(trim, "/", 3)
			if len(items) < 2 {
				return nil
			}
			cateID := items[1]

			bigCate := el.ChildText(".index-title")
			if bigCate != "美食" {
				return nil
			}

			ctx.PutReqContextValue("big_category", bigCate)
			return ctx.VisitForNextWithContext(el.Request.URL.String() + "/" + cateID)
		},
	},
}

var step3 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Infof("Visiting %s", req.URL.String())
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`#classfy a`: func(ctx *spider.Context, el *spider.HTMLElement) error {
			link := el.Attr("href")
			if link == "javascript:;" {
				return nil
			}
			subCate := el.Text
			ctx.PutReqContextValue("sub_category", subCate)

			return ctx.VisitForNextWithContext(link)
		},
	},
}

var step4 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Infof("Visiting %s", req.URL.String())
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`.navigation`: func(ctx *spider.Context, el *spider.HTMLElement) error {
			el.ForEach(`#classfy-sub a`, func(i int, element *spider.HTMLElement) {
				if i == 0 { // 第一个链接为"不限", 忽略
					return
				}
				ctx.VisitForNextWithContext(element.Attr("href"))
			})
			return nil
		},
	},
}

var step5 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Infof("Visiting %s", req.URL.String())
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`#region-nav a`: func(ctx *spider.Context, el *spider.HTMLElement) error {
			adname := el.Text
			ctx.PutReqContextValue("adname", adname)

			return ctx.VisitForNextWithContext(el.Attr("href"))
		},
	},
}

var step6 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Infof("Visiting %s", req.URL.String())
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`.page a:nth-last-child(2)`: func(ctx *spider.Context, el *spider.HTMLElement) error {
			countTxt := el.Text
			if countTxt == "" {
				return nil
			}

			count64, err := strconv.ParseInt(countTxt, 10, 64)
			if err != nil {
				log.Errorf("pase page count err:%s", err.Error())
				return nil
			}

			for i := 2; i <= int(count64); i++ {
				nextURL := fmt.Sprintf("%sp%d", el.Request.URL.String(), i)
				log.Infof("nextURL:%s", nextURL)
				ctx.Visit(nextURL)
			}
			return nil
		},
		`.shop-list li .pic a`: func(ctx *spider.Context, el *spider.HTMLElement) error {
			photo := el.ChildAttr(`img`, "src")
			ctx.PutReqContextValue("photos", photo)

			return ctx.VisitForNextWithContext(el.Attr("href"))
		},
	},
}

var step7 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Infof("Visiting %s", req.URL.String())
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`#basic-info`: func(ctx *spider.Context, el *spider.HTMLElement) error {
			shopNameNodes := el.DOM.Find(`.shop-name`).Nodes
			if len(shopNameNodes) == 0 {
				return nil
			}
			shopName := shopNameNodes[0].FirstChild.Data
			address := el.ChildText(`.address .item`)
			tel := el.ChildText(`.tel .item`)

			city := ctx.GetReqContextValue("city")
			adname := ctx.GetReqContextValue("adname")
			bigCate := ctx.GetReqContextValue("big_category")
			subCate := ctx.GetReqContextValue("sub_category")
			photos := ctx.GetReqContextValue("photos")

			return ctx.Output(map[int]interface{}{
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
}

func Retry(ctx *spider.Context, count int) error {
	req := ctx.GetRequest()
	key := fmt.Sprintf("err_req_%s", req.URL.String())

	var et int
	if errCount := ctx.GetAnyReqContextValue(key); errCount != nil {
		et = errCount.(int)
		if et >= count {
			return fmt.Errorf("exceed %d counts", count)
		}
	}
	log.Infof("errCount:%d, we wil retry url:%s, after 1 second", et+1, req.URL.String())
	time.Sleep(time.Second)
	ctx.PutReqContextValue(key, et+1)
	ctx.Retry()

	return nil
}
