package model

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/panzhenyu12/flower/entities"
	"github.com/panzhenyu12/flower/utils"
	"github.com/pkg/errors"
)

const (
	DEFAULT_USER_ID     = "1"
	DEFAULT_USER_ID_INT = 1
)

type UserModel struct {
	Ts           int64
	UserId       string `json:"Id" valid:"-"`
	UserName     string `valid:"required"`
	UserPasswd   string `valid:"required"`
	OrgID        string `json:"OrgId" valid:"required"`
	OrgName      string
	UserFuncRole *FuncRoleModel
	IsValid      bool
	RealName     string
	Comment      string
	Org          *OrgModel `json:"Org,omitempty"`
	// Status       int
	// SecurityToken string `json:"-"` // json ignore
}

func (model *UserModel) FromEntity(entity *entities.AccountExtend) {
	model.Ts = entity.Ts
	model.UserId = utils.FromInt64ToString(entity.UserId)
	model.UserName = entity.UserName
	model.OrgID = utils.FromInt64ToString(entity.OrgId)
	model.IsValid = entity.IsValid
	model.RealName = entity.RealName
	model.Comment = entity.Comment
	funcRole := &FuncRoleModel{}
	funcRole.FromEntity(&entity.FuncRole)
	model.UserFuncRole = funcRole
	if entity.Org.OrgId != 0 {
		org := new(OrgModel)
		org.FromEntity(&entity.Org)
		model.Org = org
		model.OrgName = org.OrgName
	}
}

func (model *UserModel) ToEntity() *entities.Account {
	if model == nil {
		return nil
	}
	return &entities.Account{
		UserId:     utils.FromStringToInt64(model.UserId),
		UserName:   model.UserName,
		UserPasswd: model.UserPasswd,
		OrgId:      utils.FromStringToInt64(model.OrgID),
		FuncRoleId: utils.FromStringToInt64(model.UserFuncRole.FuncRoleId),
		IsValid:    model.IsValid,
		RealName:   model.RealName,
		Comment:    model.Comment,
	}
}

type ModifyUserModel struct {
	UserId       string         `json:"Id" valid:"-"`
	UserName     string         `valid:"required"`
	UserFuncRole *FuncRoleModel `valid:"-"`
	IsValid      bool
	RealName     string
	Comment      string
	// Status       int
	// SecurityToken string `json:"-"` // json ignore
}

func (model *ModifyUserModel) Valid() error {
	if model.UserId == DEFAULT_USER_ID {
		return fmt.Errorf("Invalid id %v", model.UserId)
	}
	if model.UserFuncRole == nil || model.UserFuncRole.FuncRoleId == "" {
		return fmt.Errorf("Func role id not found")
	}
	_, err := govalidator.ValidateStruct(model)
	return errors.WithStack(err)
}

func (model *ModifyUserModel) ToEntity() *entities.Account {
	if model == nil {
		return nil
	}
	return &entities.Account{
		UserId:     utils.FromStringToInt64(model.UserId),
		UserName:   model.UserName,
		FuncRoleId: utils.FromStringToInt64(model.UserFuncRole.FuncRoleId),
		IsValid:    model.IsValid,
		RealName:   model.RealName,
		Comment:    model.Comment,
	}
}
