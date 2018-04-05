package main

import (
	"github.com/nange/gospider"
	_ "github.com/nange/gospider/_example/rule/baidunews"
	_ "github.com/nange/gospider/_example/rule/mojitianqi"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: false})
}

func main() {
	gospider.Run()
}
