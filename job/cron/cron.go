package cron

import (
	"context"
	"fmt"
	"strconv"
	"time"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	opentracing "github.com/opentracing/opentracing-go"
	opentracing_log "github.com/opentracing/opentracing-go/log"
	"github.com/robfig/cron"
)

type testJob struct {
	server *machinery.Server
}

func (job *testJob) Run() {
	// server, err := utils.StartServer()
	// if err != nil {
	// 	glog.Error(err)
	// 	return
	// }
	successTask := tasks.Signature{
		Name: "split",
	}
	concatTask := tasks.Signature{
		Name: "bang1",
		Args: []tasks.Arg{
			{
				Type:  "[]string",
				Value: []string{"foo", "bar"},
			},
		},
		OnSuccess: []*tasks.Signature{&successTask},
	}
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()
	batchID := uuid.New().String()
	span.SetBaggageItem("batch.id", batchID)
	span.LogFields(opentracing_log.String("batch.id", batchID))
	log.INFO.Println("Starting batch:", batchID)
	for i := 0; i < 2; i++ {
		//job.server.SendTask()
		concatTask.Args[0].Value = []string{"foo", strconv.Itoa(i)}
		asyncResult, err := job.server.SendTaskWithContext(ctx, &concatTask)
		if err != nil {
			fmt.Errorf("Could not send task: %s", err.Error())
			return
		}

		_, err = asyncResult.Get(time.Duration(time.Millisecond * 5))
		if err != nil {
			fmt.Errorf("Getting task result failed with error: %s", err.Error())
			return
		}
	}
}
func Startcron() {
	go StartStatsTask()
}

func StartStatsTask() {
	c := cron.New()
	c.Start()
}
