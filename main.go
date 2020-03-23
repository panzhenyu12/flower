package main

import (
	"flag"
	"fmt"
	"os"

	"flower/config"
	"flower/job"
	"flower/job/cron"
	"flower/utils"
	"flower/web/controllers"
	"flower/web/routers"

	"github.com/golang/glog"
)

var (
	debug bool = true
)

var docron = flag.Bool("cron", false, "start cron")
var doworker = flag.Bool("worker", false, "start worker")
var doweb = flag.Bool("web", debug, "start web")

func main() {
	//flag.Parse()
	conf := config.GetConfig()
	fmt.Println(conf)
	if *doweb {
		go func() {
			router := routers.NewRouter(controllers.New(conf))
			router.AddBaseRouter()
			router.RouterEngine.Run(conf.HttpServiceAddr)
		}()
		//consumer, err := kafka.NewKafkaConsumer(conf.KafkaAddress, conf.KafkaTopics)
		// if err != nil {
		// 	glog.Fatalln(err)
		// 	os.Exit(1)
		// }
	}
	if *docron {
		cron.Startcron()
	}
	if *doworker {
		go func() {
			if err := job.Worker(); err != nil {
				glog.Fatalln(err)
				os.Exit(1)
			}
		}()
	}

	sig := utils.WaitForExit() //signal exit: Ctrl+C or ...
	glog.Infof("got signal: %v, trigger to stop system", sig)
}
