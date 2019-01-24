package controllers

import (
	"net/http"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/panzhenyu12/flower/config"
	"github.com/panzhenyu12/flower/job"
	"github.com/panzhenyu12/flower/utils"
)

type PicStruct struct {
	ImageUrl string
	BinData  string
}

type BaseResp struct {
	Code     int
	ErrorMsg string
}

type SearchResp struct {
	*BaseResp
	Count int
	Data  interface{}
}

type Controller struct {
	//engine      *xorm.Engine
	redisclient *redis.Client
	jobserver   *machinery.Server
}

func New(config *config.Config) *Controller {
	return &Controller{
		redisclient: utils.GetRedisClient(),
		jobserver:   job.GetServer(),
	}
}

func (this *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (this *Controller) GetID(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"id": utils.GetIDGenerate().GetID().Int64(),
	})
}
