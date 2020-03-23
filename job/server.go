package job

import (
	"sync"

	"github.com/RichardKnop/machinery/example/tracers"
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/golang/glog"
	localconfig "flower/config"
	"flower/job/jobs"
	"github.com/pkg/errors"
)

var slock sync.Mutex
var jobserver *machinery.Server

var tasks = map[string]interface{}{
	"add":      jobs.Add,
	"multiply": jobs.Multiply,
}

func loadConfig() (*config.Config, error) {
	configPath := localconfig.GetConfig().TaskConfig
	if configPath != "" {
		return config.NewFromYaml(configPath, true)
	}
	return config.NewFromEnvironment(true)
}
func GetServer() *machinery.Server {
	if jobserver == nil {
		slock.Lock()
		defer slock.Unlock()
		if jobserver == nil {
			cnf, err := loadConfig()
			glog.Info("init jobserver")
			if err != nil {
				glog.Error(err)
				return nil
			}
			// Create server instance
			server, err := machinery.NewServer(cnf)
			if err != nil {
				glog.Error(err)
				return nil
			}
			jobserver = server
		}
	}
	return jobserver
}
func StartServer() (*machinery.Server, error) {
	server := GetServer()
	if server == nil {
		return nil, errors.New("init job server error")
	}
	return server, server.RegisterTasks(tasks)
}

func Worker() error {
	//ConsumerTag
	consumerTag := localconfig.GetConfig().ConsumerTag
	cleanup, err := tracers.SetupTracer(consumerTag)
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()
	server, err := StartServer()
	if err != nil {
		return errors.WithStack(err)
	}
	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker(consumerTag, 100)
	return worker.Launch()
}
func initTasks() {
	tasks = map[string]interface{}{
		"add":               jobs.Add,
		"multiply":          jobs.Multiply,
		"sum_ints":          jobs.SumInts,
		"sum_floats":        jobs.SumFloats,
		"concat":            jobs.Concat,
		"split":             jobs.Split,
		"panic_task":        jobs.PanicTask,
		"long_running_task": jobs.LongRunningTask,
	}
}
func RegisterTasks(fname string, job interface{}) {
	if tasks == nil || len(tasks) == 0 {
		initTasks()
	}
	tasks[fname] = job
}
