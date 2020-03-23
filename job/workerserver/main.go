package main

import (
	"os"

	"github.com/golang/glog"
	"flower/job"
)

func main() {
	if err := job.Worker(); err != nil {
		glog.Fatalln(err)
		os.Exit(1)
	}
}
