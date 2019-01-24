package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/panzhenyu12/flower/common"
	"github.com/panzhenyu12/flower/model"
	"github.com/panzhenyu12/flower/repositories"
	"github.com/panzhenyu12/flower/utils"
	"github.com/pkg/errors"
)

func (controller *Controller) QueryFuncRoles(ctx *gin.Context) {
	query := new(model.FuncRoleQuery)
	if err := ShouldBindValidJSON(ctx, query); err != nil {
		Error400(ctx, err)
		return
	}
	dal := repositories.GetFuncRoleDal()
	entities, count, err := dal.QueryAndCount(query)
	if err != nil {
		Error500(ctx, err)
		return
	}
	models := make([]*model.FuncRoleModel, len(entities))
	for i, e := range entities {
		m := &model.FuncRoleModel{}
		m.FromEntity(e)
		models[i] = m
	}
	SearchRespJSON(ctx, models, count)
}

func (controller *Controller) AddFuncRole(ctx *gin.Context) {
	m := new(model.FuncRoleModel)
	if err := ShouldBindValidJSON(ctx, m); err != nil {
		Error400(ctx, err)
		return
	}
	e := m.ToEntity()
	if e.FuncRoleId <= 0 {
		id := utils.GetIDGenerate().GetID()
		e.FuncRoleId = id.Int64()
		e.Ts = id.Time()
		e.Uts = time.Now()
		e.Status = int(common.TableStatus_Create)
	}
	dal := repositories.GetFuncRoleDal()
	if err := dal.Insert(e, nil); err != nil {
		Error500(ctx, err)
		return
	}
	NoContent(ctx)
}

func (controller *Controller) ModifyFuncRole(ctx *gin.Context) {
	m := new(model.ModifyFuncRoleModel)
	if err := ShouldBindValidJSON(ctx, m); err != nil {
		Error400(ctx, err)
		return
	}
	//只允许并发一次修改
	fmt.Println(m.FuncRoleId)
	rlock, err := utils.Obtain(controller.redisclient, m.FuncRoleId, nil)
	if err != nil {
		Error500(ctx, err)
		return
	}
	defer rlock.Unlock()
	e := m.ToEntity()
	dal := repositories.GetFuncRoleDal()
	if _, err := dal.Transaction(func(session *repositories.Session) (interface{}, error) {
		if err := dal.Upate(e, session); err != nil {
			return nil, errors.WithStack(err)
		}
		userIds, err := dal.GetUserIds(session, e.FuncRoleId)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if len(userIds) == 0 {
			return nil, nil
		}
		m := make(map[int64]string)
		for _, userId := range userIds {
			m[userId] = utils.GetIDGenerate().GetID().String()
		}
		accountdal := repositories.GetAccountDal()
		if err := accountdal.BatchUpdateSecurityToken(m, session); err != nil {
			return nil, errors.WithStack(err)
		}
		return nil, nil
	}); err != nil {
		Error500(ctx, err)
		return
	}
	NoContent(ctx)
}

func (controller *Controller) DeleteFuncRole(ctx *gin.Context) {
	id, err := ParamInt64(ctx, "id")
	if err != nil {
		Error400(ctx, err)
		return
	} else if id == model.DEFAULT_FUNC_ROLE_ID_INT {
		Error400f(ctx, "Invalid id %v", id)
		return
	}
	if err := controller.deleteFuncRole(id); err != nil {
		Error500(ctx, err)
		return
	}
	NoContent(ctx)
}

func (controller *Controller) DeleteFuncRoles(ctx *gin.Context) {
	var ids []string
	if err := ShouldBindValidJSON(ctx, &ids); err != nil {
		Error400(ctx, err)
		return
	} else if len(ids) == 0 {
		Error400f(ctx, "Ids not found")
		return
	}
	result := make(map[string]string)
	for _, id := range ids {
		if id == model.DEFAULT_FUNC_ROLE_ID {
			msg := fmt.Sprintf("Invalid id %v", id)
			glog.Errorln(msg)
			result[id] = msg
			continue
		}
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			msg := fmt.Sprintf("Invalid id %v", id)
			glog.Errorln(msg)
			result[id] = msg
			continue
		}
		if err := controller.deleteFuncRole(intId); err != nil {
			glog.Errorln(err)
			result[id] = err.Error()
			continue
		}
		result[id] = ""
	}
	OKJSON(ctx, result)
}

func (controller *Controller) deleteFuncRole(id int64) error {
	//只允许并发一次修改
	rlock, err := utils.Obtain(controller.redisclient, strconv.FormatInt(id, 10), nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer rlock.Unlock()
	dal := repositories.GetFuncRoleDal()
	if _, err := dal.Transaction(func(session *repositories.Session) (interface{}, error) {
		userIds, err := dal.GetUserIds(session, id)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if len(userIds) > 0 {
			return nil, errors.New("Entity in use") // todo code
		}
		err = dal.SoftDelete(id, session)
		return nil, errors.WithStack(err)
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
