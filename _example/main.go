package main

import (
	"github.com/nange/gospider"
	_ "github.com/nange/gospider/_example/rule/baidunews"
	_ "github.com/nange/gospider/_example/rule/dianping"
	_ "github.com/nange/gospider/_example/rule/mojitianqi"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"})
}

func main() {
	gs := gospider.New()
	logrus.Fatal(gs.Run())
}
