package model

import (
	"fmt"

	"github.com/panzhenyu12/flower/utils"

	"github.com/asaskevich/govalidator"
	"github.com/panzhenyu12/flower/entities"
	"github.com/pkg/errors"
)

const (
	DEFAULT_FUNC_ROLE_ID     = "1"
	DEFAULT_FUNC_ROLE_ID_INT = 1
)

type FuncRoleModel struct {
	Ts           int64  `valid:"-"`
	FuncRoleId   string `json:"Id" valid:"-"`
	FuncRoleName string `json:"Name" valid:"required"`
	Content      string `valid:"required"`
	Comment      string `valid:"-"`
	UserCount    int    `valid:"-"`
}

func (model *FuncRoleModel) FromEntity(e *entities.FuncRole) {
	if e == nil {
		return
	}
	model.Ts = e.Ts
	model.FuncRoleId = utils.FromInt64ToString(e.FuncRoleId)
	model.FuncRoleName = e.FuncRoleName
	model.Content = e.Content
	model.Comment = e.Comment
	model.UserCount = e.UserCount
}

func (model *FuncRoleModel) ToEntity() *entities.FuncRole {
	if model == nil {
		return nil
	}
	return &entities.FuncRole{
		Ts:           model.Ts,
		FuncRoleId:   utils.FromStringToInt64(model.FuncRoleId),
		FuncRoleName: model.FuncRoleName,
		Content:      model.Content,
		Comment:      model.Comment,
	}
}

type ModifyFuncRoleModel struct {
	FuncRoleId   string `json:"Id" valid:"required"`
	FuncRoleName string `json:"Name" valid:"required"`
	Content      string `valid:"required"`
	Comment      string `valid:"-"`
}

func (model *ModifyFuncRoleModel) Valid() error {
	if model.FuncRoleId == DEFAULT_FUNC_ROLE_ID {
		return fmt.Errorf("Invalid id %v", model.FuncRoleId)
	}
	_, err := govalidator.ValidateStruct(model)
	return errors.WithStack(err)
}

func (model *ModifyFuncRoleModel) ToEntity() *entities.FuncRole {
	if model == nil {
		return nil
	}
	return &entities.FuncRole{
		FuncRoleId:   utils.FromStringToInt64(model.FuncRoleId),
		FuncRoleName: model.FuncRoleName,
		Content:      model.Content,
		Comment:      model.Comment,
	}
}
