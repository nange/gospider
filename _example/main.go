package main

import (
	"github.com/nange/gospider"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	gospider.Run()
}
