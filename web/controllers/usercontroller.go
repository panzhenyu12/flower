package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/panzhenyu12/flower/entities"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/panzhenyu12/flower/common"
	"github.com/panzhenyu12/flower/model"
	"github.com/panzhenyu12/flower/repositories"
	"github.com/panzhenyu12/flower/utils"
)

func (controller *Controller) GetCurrentUser(ctx *gin.Context) {
	account := GetCurrentUserFromCtx(ctx)
	if account == nil {
		Error404f(ctx, "User not found")
	}
	user := &model.UserModel{}
	user.FromEntity(account)
	OKJSON(ctx, user)
}

func (controller *Controller) AddUser(ctx *gin.Context) {
	m := new(model.UserModel)
	if err := ShouldBindValidJSON(ctx, m); err != nil {
		Error400(ctx, err)
		return
	}
	e := m.ToEntity()
	if e.UserId <= 0 {
		id := utils.GetIDGenerate().GetID()
		e.UserId = id.Int64()
		e.Ts = id.Time()
		e.Uts = time.Now()
		e.Status = int(common.TableStatus_Create)
	}
	e.SecurityToken = genSecurityToken()
	dal := repositories.GetAccountDal()
	if err := dal.Insert(e); err != nil {
		Error500(ctx, err)
		return
	}
	NoContent(ctx)
}

func (controller *Controller) ModifyUser(ctx *gin.Context) {
	m := new(model.ModifyUserModel)
	if err := ShouldBindValidJSON(ctx, m); err != nil {
		Error400(ctx, err)
		return
	}
	//只允许并发一次修改
	rlock, err := utils.Obtain(controller.redisclient, m.UserId, nil)
	if err != nil {
		Error500(ctx, err)
		return
	}
	defer rlock.Unlock()
	e := m.ToEntity()
	e.SecurityToken = genSecurityToken()
	dal := repositories.GetAccountDal()
	if err := dal.Update(e); err != nil {
		Error500(ctx, err)
		return
	}
	NoContent(ctx)
}

func (controller *Controller) DeleteUser(ctx *gin.Context) {
	id, err := ParamInt64(ctx, "id")
	if err != nil {
		Error400(ctx, err)
		return
	} else if id == model.DEFAULT_USER_ID_INT {
		Error400f(ctx, "Invalid id %v", id)
		return
	}
	if err := controller.deleteUser(id); err != nil {
		Error500(ctx, err)
		return
	}
	NoContent(ctx)
}

func (controller *Controller) DeleteUsers(ctx *gin.Context) {
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
		if id == model.DEFAULT_USER_ID {
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
		if err := controller.deleteUser(intId); err != nil {
			glog.Errorln(err)
			result[id] = err.Error()
			continue
		}
		result[id] = ""
	}
	OKJSON(ctx, result)
}

func (controller *Controller) ModifyUserValid(ctx *gin.Context) {
	controller.modifyUserIsValid(ctx, true)
}

func (controller *Controller) ModifyUserInvalid(ctx *gin.Context) {
	controller.modifyUserIsValid(ctx, false)
}

func (controller *Controller) QueryUsers(ctx *gin.Context) {
	query := new(model.UsersQuery)
	err := ShouldBindValidJSON(ctx, query)
	if err != nil {
		Error400(ctx, err)
		return
	}
	id := utils.FromStringToInt64(query.OrgID)
	if id == 0 {
		Error400(ctx, errors.New("OrgID error"))
		return
	}
	downorgs, err := repositories.GetDeepOrgDal().GetDownOrg(id)
	if err != nil {
		Error500(ctx, err)
		return
	}
	orgids := make([]int64, 0)
	for _, org := range downorgs {
		orgids = append(orgids, org.OrgId)
	}
	count, data, err := repositories.GetAccountDal().QueryAndCount(query.BaseQuery, orgids)
	if err != nil {
		Error500(ctx, err)
		return
	}
	users := make([]*model.UserModel, 0)
	for _, d := range data {
		u := new(model.UserModel)
		u.FromEntity(d)
		users = append(users, u)
	}
	ret := model.ResponseData{
		Count: count,
		Data:  users,
	}
	OKJSON(ctx, ret)
}

func (controller *Controller) UpdatePassword(ctx *gin.Context) {
	var req model.UpdatePasswordRequest
	if err := ShouldBindValidJSON(ctx, &req); err != nil {
		Error400(ctx, err)
		return
	}
	user := GetCurrentUserFromCtx(ctx)
	if user == nil {
		Error404f(ctx, "User not found")
		return
	}
	//只允许并发一次修改
	rlock, err := utils.Obtain(controller.redisclient, fmt.Sprintf("%v", user.UserId), nil)
	if err != nil {
		Error500(ctx, err)
		return
	}
	defer rlock.Unlock()
	dal := repositories.GetAccountDal()
	e, err := dal.GetValidByUserName(user.UserName)
	if err != nil {
		Error500(ctx, err)
		return
	}
	if e.UserPasswd != model.EncryptPassword(req.OldPassword) {
		Error400f(ctx, "Password not match")
		return
	}
	e = &entities.Account{
		UserId:        user.UserId,
		UserPasswd:    model.EncryptPassword(req.NewPassword),
		SecurityToken: utils.GetIDGenerate().GetID().String(),
	}
	if err := dal.UpdateCore(e, nil); err != nil {
		Error500(ctx, err)
		return
	}
	NoContent(ctx)
}

func (controller *Controller) deleteUser(id int64) error {
	//只允许并发一次修改
	rlock, err := utils.Obtain(controller.redisclient, strconv.FormatInt(id, 10), nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer rlock.Unlock()
	dal := repositories.GetAccountDal()
	if err := dal.Delete(id); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (controller *Controller) modifyUserIsValid(ctx *gin.Context, isValid bool) {
	id, err := ParamInt64(ctx, "id")
	if err != nil {
		Error400(ctx, err)
		return
	} else if id == model.DEFAULT_USER_ID_INT {
		Error400f(ctx, "Invalid id %v", id)
		return
	}
	//只允许并发一次修改
	rlock, err := utils.Obtain(controller.redisclient, ctx.Param("id"), nil)
	if err != nil {
		Error500(ctx, err)
		return
	}
	defer rlock.Unlock()
	e := &entities.Account{
		UserId:        id,
		IsValid:       isValid,
		SecurityToken: genSecurityToken(),
	}
	dal := repositories.GetAccountDal()
	s := dal.GetNewSession()
	s.MustCols("is_valid")
	if err := dal.UpdateCore(e, s); err != nil {
		Error500(ctx, err)
		return
	}
	NoContent(ctx)
}

func genSecurityToken() string {
	return utils.GetIDGenerate().GetID().String()
}
