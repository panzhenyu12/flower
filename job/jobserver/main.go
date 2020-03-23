package main

import (
	"fmt"

	"github.com/golang/glog"
	"flower/config"
	"flower/job/cron"
	"flower/utils"
)

func main() {
	conf := config.GetConfig()
	fmt.Println(conf)
	cron.Startcron()
	sig := utils.WaitForExit() //signal exit: Ctrl+C or ...
	glog.Infof("got signal: %v, trigger to stop system", sig)
}
