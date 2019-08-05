package rule

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/nange/gospider/spider"
)

type GetRuleListResp struct {
	Total int      `json:"total"`
	Data  []string `json:"data"`
}

func GetRuleList(c *gin.Context) {
	keys := spider.GetTaskRuleKeys()
	if len(keys) == 0 {
		log.Warnf("task rule is empty")
	} else {
		sort.Sort(Pinyin(keys))
	}

	c.JSON(http.StatusOK, &GetRuleListResp{
		Total: len(keys),
		Data:  keys,
	})
}

type Pinyin []string

func (s Pinyin) Len() int      { return len(s) }
func (s Pinyin) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Pinyin) Less(i, j int) bool {
	a, _ := UTF82GB18030(s[i])
	b, _ := UTF82GB18030(s[j])
	bLen := len(b)
	for idx, chr := range a {
		if idx > bLen-1 {
			return false
		}
		if chr != b[idx] {
			return chr < b[idx]
		}
	}
	return true
}

func UTF82GB18030(src string) ([]byte, error) {
	GB18030 := simplifiedchinese.All[0]
	return ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(src)), GB18030.NewEncoder()))
}
