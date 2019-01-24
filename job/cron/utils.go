package cron

import (
	"time"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends/result"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/golang/glog"
)

func dojob(category string, server *machinery.Server, signatures ...*tasks.Signature) {
	group, err := tasks.NewGroup(signatures...)
	if err != nil {
		glog.Error(category, err)
		return
	}
	asyncResults, err := server.SendGroup(group, 0)
	for _, asyncResult := range asyncResults {
		go func(asyncResult *result.AsyncResult) {
			_, err := asyncResult.Get(time.Duration(time.Millisecond * 1))
			if err != nil {
				glog.Error(category, err, asyncResult.Signature)
			}
		}(asyncResult)
	}
}
