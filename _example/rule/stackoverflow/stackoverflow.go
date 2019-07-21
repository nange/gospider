package stackoverflow

import (
	"github.com/nange/gospider/spider"
	log "github.com/sirupsen/logrus"
)

func init() {
	spider.Register(rule)
}

var (
	outputFields = []string{"question_title", "question_detail", "question_tags", "answer_list"}
	constraints  = spider.NewConstraints(outputFields,
		"VARCHAR(512) NOT NULL DEFAULT ''",
		"TEXT",
		"VARCHAR(512) NOT NULL DEFAULT ''",
		"TEXT",
	)
)
var rule = &spider.TaskRule{
	Name:              "StackOverFlow",
	Description:       "StackOverFlow Highly Quality QA",
	Namespace:         "stackoverflow_en",
	OutputFields:      outputFields,
	OutputConstraints: constraints,
	Rule: &spider.Rule{
		Head: func(ctx *spider.Context) error {
			return ctx.VisitForNext("https://stackoverflow.com/questions?tab=votes&page=1")
		},
		Nodes: map[int]*spider.Node{
			0: step1,
			1: step2,
		},
	},
}
var step1 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Infof("Visiting %s", req.URL.String())
	},
	OnError: func(ctx *spider.Context, res *spider.Response, err error) error {
		log.Errorf("Visiting failed! url:%s,err:%s", res.Request.URL.String(), err.Error())
		return Retry(ctx, 3)
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`.summary h3 a`: func(context *spider.Context, element *spider.HTMLElement) error {
			link := element.Attr("href")
			link = "https://stackoverflow.com" + link
			return context.VisitForNext(link)
		},
	},
}
var step2 = &spider.Node{
	OnRequest: func(ctx *spider.Context, req *spider.Request) {
		log.Println("Visting", req.URL.String())
	},
	OnHTML: map[string]func(*spider.Context, *spider.HTMLElement) error{
		`.inner-content.clearfix`: func(ctx *spider.Context, element *spider.HTMLElement) error {
			question_title := element.ChildText("#question-header h1 a")

			question_detail, err := element.DOM.Find("#mainbar .question .post-layout .postcell.post-layout--right .post-text").Html()
			if err != nil {
				log.Errorf("step2 question_detail canot find the quesion:", question_title, err.Error())
				question_detail = "No Descrption"
			}

			question_taglist := ""
			element.ForEach("#mainbar .question .post-layout .postcell.post-layout--right .post-taglist.grid.gs4.gsy.fd-column .grid.ps-relative.d-block a",
				func(i int, tagEle *spider.HTMLElement) {
					tag := tagEle.Attr("href")
					if tag == "" {
						log.Errorf("step2 question_tag canot find  the question:", question_title)
					} else {
						question_taglist += (tag + "T^T")
					}
				})

			answer_detail := ""
			accept_answer := element.DOM.Find("#mainbar #answers .answer.accepted-answer")
			if accept_answer != nil {
				accept_answer = accept_answer.Find(".post-layout .answercell.post-layout--right .post-text")
				if accept_answer != nil {
					answer_detail, err = accept_answer.Html()
					if err != nil {
						log.Errorf("step2 acceptAnswer.Html() error ,quesionTitle is:", question_title, err.Error())
					}
				} else {
					//todo: error about accept answer find no text
					log.Errorf("mainbar #answers .answer.accepted-answer find but .post-layout .answercell.post-layout--right .post-text not find ; question title", question_title)
				}
			} else {
				//no accept answer: find first answer
				accept_answer = accept_answer.Find(".answer").First()
				if accept_answer != nil {
					accept_answer = accept_answer.Find(".post-layout .answercell.post-layout--right .post-text")
					if accept_answer != nil {
						answer_detail, err = accept_answer.Html()
						if err != nil {
							log.Errorf("step2 Answer.Html() error ,quesionTitle is:", question_title, err.Error())
						}
					} else {
						log.Errorf("step2 canot find  post-layout .answercell.post-layout--right .post-text ;quesionTitle is:", question_title)
					}
				} else {
					log.Errorf("step2 canot find  anser first ;quesionTitle is:", question_title)
				}
			}
			return ctx.Output(map[int]interface{}{
				0: question_title,
				1: question_detail,
				2: question_taglist,
				3: answer_detail,
			})
		},
	},
}

func Retry(ctx *spider.Context, count int) error {
	log.Errorf("need to retry")
	return nil
}
