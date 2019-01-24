package main

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/panzhenyu12/flower/config"
	"github.com/panzhenyu12/flower/job/cron"
	"github.com/panzhenyu12/flower/utils"
)

func main() {
	conf := config.GetConfig()
	fmt.Println(conf)
	cron.Startcron()
	sig := utils.WaitForExit() //signal exit: Ctrl+C or ...
	glog.Infof("got signal: %v, trigger to stop system", sig)
}
