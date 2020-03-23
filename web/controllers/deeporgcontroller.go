package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"flower/model"
	"flower/repositories"
	"flower/utils"
	"github.com/pkg/errors"
)

func (controller *Controller) GetOrg(ctx *gin.Context) {
	var id int64
	if str := ctx.Param("id"); str != "" {
		if val, err := strconv.ParseInt(str, 10, 64); err != nil {
			Error400f(ctx, "Invalid id %v", str)
			return
		} else {
			id = val
		}
	} else {
		id = 1
		// todo id = current user's org id
	}
	var query model.OrgQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		Error400(ctx, err)
		return
	}
	// query.OrgId = id
	orgdal := repositories.GetDeepOrgDal()
	orgs, err := orgdal.GetDownOrg(id)
	if err != nil {
		Error500(ctx, err)
		return
	}
	orgMap := make(map[string]*model.OrgModel)
	orgIds := make([]int64, len(orgs))
	var root *model.OrgModel
	var maxLevel int64
	maxLevel = -1
	for i, e := range orgs {
		orgIds[i] = e.OrgId
		m := &model.OrgModel{}
		m.FromEntity(e)
		orgMap[m.OrgID] = m
		if e.OrgId == id {
			root = m
			if query.RecursiveDepth >= 0 {
				maxLevel = m.OrgLevel + query.RecursiveDepth
			}
			continue
		}
		if maxLevel < 0 || m.OrgLevel <= maxLevel {
			// DO NOT BREAK!!!!
			// orgMap IS NEEDED
			if parent, exist := orgMap[m.SuperiorOrgID]; exist {
				parent.Child = append(parent.Child, m)
			}
		}
	}

	if query.IncludeUsers {
		accountdal := repositories.GetAccountDal()
		accounts, err := accountdal.GetByOrgIDs(orgIds)
		if err != nil {
			Error500(ctx, err)
			return
		}
		for _, e := range accounts {
			m := &model.UserModel{}
			m.FromEntity(e)
			org := orgMap[m.OrgID]
			m.OrgName = org.OrgName
			org.Users = append(org.Users, m)
		}
	}

	OKJSON(ctx, root)
}

func (controller *Controller) AddOrg(ctx *gin.Context) {
	org := new(model.OrgModel)
	resp := &model.BaseResponse{}
	err := ctx.ShouldBindJSON(org)
	if err != nil {
		resp.ErrorMsg = err.Error()
		glog.Error(err)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	err = org.Valid()
	if err != nil {
		glog.Error(err)
		resp.ErrorMsg = err.Error()
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	orgdal := repositories.GetDeepOrgDal()
	if utils.FromStringToInt64(org.SuperiorOrgID) > 0 {
		//获取父级组织
		superiororg, err := orgdal.GetOrgByID(utils.FromStringToInt64(org.SuperiorOrgID))
		if err != nil {
			//msg := "get "
			glog.Error(err)
			resp.ErrorMsg = err.Error()
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
		if superiororg == nil || superiororg.OrgId == 0 {
			msg := "get superiororg error"
			glog.Error(msg)
			resp.ErrorMsg = msg
			ctx.JSON(http.StatusBadRequest, resp)
			return
		}
		org.OrgLevel = superiororg.OrgLevel + 1
	} else {
		org.OrgLevel = 1
	}
	orgentity := org.ToEntity()
	err = orgdal.AddOne(orgentity, nil)
	if err != nil {
		glog.Error(err)
		resp.ErrorMsg = err.Error()
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	ctx.JSON(http.StatusOK, orgentity)
}

func (controller *Controller) ModifyOrg(ctx *gin.Context) {
	orgreq := new(model.ModifyOrgModel)
	resp := &model.BaseResponse{}
	err := ctx.ShouldBindJSON(orgreq)
	if err != nil {
		resp.ErrorMsg = err.Error()
		glog.Error(err)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	err = orgreq.Valid()
	if err != nil {
		resp.ErrorMsg = err.Error()
		glog.Error(err)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	//只允许并发一次修改
	rlock, err := utils.Obtain(controller.redisclient, orgreq.OrgID, nil)
	if err != nil {
		resp.ErrorMsg = err.Error()
		glog.Error(err)
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	defer rlock.Unlock()
	orgdal := repositories.GetDeepOrgDal()
	entity, err := orgdal.GetOrgByID(utils.FromStringToInt64(orgreq.OrgID))
	if err != nil {
		resp.ErrorMsg = err.Error()
		glog.Error(err)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	entity.OrgName = orgreq.OrgName
	entity.Comment = orgreq.Comment
	err = orgdal.UpDate(entity, nil)
	if err != nil {
		resp.ErrorMsg = err.Error()
		glog.Error(err)
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Data = true
	ctx.JSON(http.StatusOK, resp)
}

func (controller *Controller) DeleteOrg(ctx *gin.Context) {
	id, err := ParamInt64(ctx, "id")
	if err != nil {
		Error400(ctx, err)
		return
	} else if id == model.DEFAULT_ORG_ID_INT {
		Error400f(ctx, "Invalid id %v", id)
		return
	}
	if err := controller.deleteOrg(id); err != nil {
		Error500(ctx, err)
		return
	}
	NoContent(ctx)
}

func (controller *Controller) DeleteOrgs(ctx *gin.Context) {
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
		if id == model.DEFAULT_ORG_ID {
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
		if err := controller.deleteOrg(intId); err != nil {
			glog.Errorln(err)
			result[id] = err.Error()
			continue
		}
		result[id] = ""
	}
	OKJSON(ctx, result)
}

func (controller *Controller) deleteOrg(id int64) error {
	dal := repositories.GetDeepOrgDal()
	_, err := dal.Transaction(func(session *repositories.Session) (interface{}, error) {
		if got, err := dal.HasChild(id, session); err != nil {
			return nil, errors.WithStack(err)
		} else if got {
			return nil, errors.New("Entity in use")
		}

		if got, err := dal.HasAccount(id, session); err != nil {
			return nil, errors.WithStack(err)
		} else if got {
			return nil, errors.New("Entity in use")
		}
		return nil, errors.WithStack(dal.Delete(session, id))
	})
	return errors.WithStack(err)
}
