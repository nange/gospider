package mojitianqi

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/nange/gospider/spider"
	"github.com/sirupsen/logrus"
)

func init() {
	spider.Register(rule)
}

var rule = &spider.TaskRule{
	Name:         "墨迹天气全国空气质量",
	Description:  "抓取墨迹天气全国各个城市区县空气质量数据",
	Namespace:    "moji_tianqi",
	OutputFields: []string{"province", "area", "aqi", "quality_grade", "pm10", "pm25", "no2", "so2", "o3", "co", "tip", "publish_time"},
	Rule: &spider.Rule{
		Head: func(ctx *spider.Context) error { // 定义入口
			return ctx.VisitForNext("https://tianqi.moji.com/aqi/china")
		},
		Nodes: map[int]*spider.Node{
			0: &spider.Node{ // 第一步: 找到全国各省城市区县的链接
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Infof("Visiting %s", req.URL.String())
				},
				OnError: func(ctx *spider.Context, res *spider.Response, err error) {
					logrus.Errorf("Visiting failed! url:%s, err:%s", res.Request.URL.String(), err.Error())
					// 出错时重试三次
					Retry(res.Request, 3)
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement){
					`.city a`: func(ctx *spider.Context, el *spider.HTMLElement) {
						link := el.Attr("href")
						el.Request.Visit(link)
					},
					`.city_hot a`: func(ctx *spider.Context, el *spider.HTMLElement) {
						link := el.Attr("href")
						el.Request.VisitForNext(link)
					},
				},
			},
			1: &spider.Node{ // 第二步: 爬取各城市区县页面上具体的空气质量数据
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Infof("Visiting %s", req.URL.String())
				},
				OnError: func(ctx *spider.Context, res *spider.Response, err error) {
					logrus.Errorf("Visiting failed! url:%s, err:%s", res.Request.URL.String(), err.Error())
					// 出错时重试三次
					Retry(res.Request, 3)
				},
				OnHTML: map[string]func(*spider.Context, *spider.HTMLElement){
					`body`: func(ctx *spider.Context, body *spider.HTMLElement) {
						var pm10, pm25, no2, so2, o3, co, publishTime string
						body.ForEach(`#aqi_info li span`, func(i int, element *spider.HTMLElement) {
							ret := element.Text
							switch i {
							case 0:
								pm10 = ret
							case 1:
								pm25 = ret
							case 2:
								no2 = ret
							case 3:
								so2 = ret
							case 4:
								o3 = ret
							case 5:
								co = ret
							}
						})
						aqi := body.ChildText("#aqi_value")
						qualityGrade := body.ChildText(`#aqi_desc`)
						publishTime = body.ChildText(".aqi_info_time b")
						publishTime = strings.TrimLeft(publishTime, "发布日期：")

						body.Request.PutReqContextValue("aqi", aqi)
						body.Request.PutReqContextValue("quality_grade", qualityGrade)
						body.Request.PutReqContextValue("pm10", pm10)
						body.Request.PutReqContextValue("pm25", pm25)
						body.Request.PutReqContextValue("no2", no2)
						body.Request.PutReqContextValue("so2", so2)
						body.Request.PutReqContextValue("o3", o3)
						body.Request.PutReqContextValue("co", co)
						body.Request.PutReqContextValue("publish_time", publishTime)

						province := body.ChildText(`.crumb li:nth-last-child(2)`)
						area := body.ChildText(`.crumb li:nth-last-child(1)`)
						body.Request.PutReqContextValue("province", province)
						body.Request.PutReqContextValue("area", area)

						internalID := body.ChildAttr(`#internal_id`, "value")
						if internalID == "" {
							return
						}
						link := fmt.Sprintf("https://tianqi.moji.com/api/getAqi/%s", internalID)
						body.Request.VisitForNextWithContext(link)
					},
				},
			},
			2: &spider.Node{ // 第三步: 由于tips字段是另外单独的请求,所以第四步单独获取tips(温馨提示)字段内容
				OnRequest: func(ctx *spider.Context, req *spider.Request) {
					logrus.Infof("Visiting %s", req.URL.String())
				},
				OnError: func(ctx *spider.Context, res *spider.Response, err error) {
					logrus.Errorf("Visiting failed! url:%s, err:%s", res.Request.URL.String(), err.Error())
					// 出错时重试三次
					Retry(res.Request, 3)
				},
				OnResponse: func(ctx *spider.Context, res *spider.Response) {
					type tip struct {
						Tips string `json:"tips"`
					}
					var ret tip
					if err := json.Unmarshal(res.Body, &ret); err != nil {
						logrus.Errorf("Unmarshal tips err:%s, body:%s", err.Error(), string(res.Body))
					}
					tips := ret.Tips
					province := res.Request.GetReqContextValue("province")
					area := res.Request.GetReqContextValue("area")
					aqi := res.Request.GetReqContextValue("aqi")
					qualityGrade := res.Request.GetReqContextValue("quality_grade")
					pm10 := res.Request.GetReqContextValue("pm10")
					pm25 := res.Request.GetReqContextValue("pm25")
					no2 := res.Request.GetReqContextValue("no2")
					so2 := res.Request.GetReqContextValue("so2")
					o3 := res.Request.GetReqContextValue("o3")
					co := res.Request.GetReqContextValue("co")
					publishTime := res.Request.GetReqContextValue("publish_time")

					ctx.Output(map[int]interface{}{
						0:  province,
						1:  area,
						2:  aqi,
						3:  qualityGrade,
						4:  pm10,
						5:  pm25,
						6:  no2,
						7:  so2,
						8:  o3,
						9:  co,
						10: tips,
						11: publishTime,
					})
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
