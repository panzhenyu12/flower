package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/panzhenyu12/flower/entities"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/panzhenyu12/flower/model"
	"github.com/pkg/errors"
)

const (
	CTX_KEY_USER = "USER"
)

func GetCurrentUserFromCtx(ctx *gin.Context) *entities.AccountExtend {
	data, exist := ctx.Get(CTX_KEY_USER)
	if !exist {
		return nil
	}
	user, ok := data.(*entities.AccountExtend)
	if !ok {
		return nil
	}
	return user
}

func ShouldBindValidJSON(ctx *gin.Context, obj interface{}) error {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		return errors.WithStack(err)
	}
	if validator, ok := obj.(model.Validator); ok {
		return errors.WithStack(validator.Valid())
	}
	return nil
}

func ParamInt64(ctx *gin.Context, key string) (int64, error) {
	val := ctx.Param(key)
	if val == "" {
		return 0, fmt.Errorf("Param %v not found", key)
	}
	i, err := strconv.ParseInt(val, 10, 64)
	return i, errors.WithStack(err)
}

func GetQueryBool(ctx *gin.Context, key string) bool {
	val := ctx.DefaultQuery(key, "")
	return strToBool(val)
}

func strToBool(val string) bool {
	return val == "1" || strings.ToLower(val) == "true"
}

func SearchRespJSON(ctx *gin.Context, results interface{}, count int) {
	obj := &SearchResp{
		Data:  results,
		Count: count,
	}
	OKJSON(ctx, obj)
}

func OKJSON(ctx *gin.Context, obj interface{}) {
	ctx.JSON(http.StatusOK, obj)
}

func NoContent(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}

func Error404f(ctx *gin.Context, format string, args ...interface{}) {
	Error404(ctx, fmt.Errorf(format, args...))
}

func Error404(ctx *gin.Context, err error) {
	Error(ctx, http.StatusNotFound, err)
}

func Error400f(ctx *gin.Context, format string, args ...interface{}) {
	Error400(ctx, fmt.Errorf(format, args...))
}

func Error400(ctx *gin.Context, err error) {
	Error(ctx, http.StatusBadRequest, err)
}

func Error401f(ctx *gin.Context, format string, args ...interface{}) {
	Error401(ctx, fmt.Errorf(format, args...))
}

func Error401(ctx *gin.Context, err error) {
	Error(ctx, http.StatusUnauthorized, err)
}

func Error500f(ctx *gin.Context, format string, args ...interface{}) {
	Error500(ctx, fmt.Errorf(format, args...))
}

func Error500(ctx *gin.Context, err error) {
	Error(ctx, http.StatusInternalServerError, err)
}

func Error(ctx *gin.Context, code int, err error) {
	glog.Errorf("%+v", err)
	MsgJSON(ctx, code, err.Error())
}

func MsgJSON(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, &BaseResp{
		Code:     code,
		ErrorMsg: msg,
	})
}
