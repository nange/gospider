package spider

import (
	"log"
	"testing"
	"time"

	"github.com/robfig/cron"
)

type testJob struct {
	c *cron.Cron
}

func (j testJob) Run() {
	log.Println("run once...")
}

func TestCron(t *testing.T) {
	c := cron.New()
	c.AddJob("0 */1 * * * *", testJob{c: c})
	c.Start()

	time.Sleep(5 * time.Minute)
}
