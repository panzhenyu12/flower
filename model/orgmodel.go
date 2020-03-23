package model

import (
	"errors"
	"time"

	"flower/utils"

	"github.com/asaskevich/govalidator"
	//"flower/common"
	"flower/entities"
)

const (
	DEFAULT_ORG_ID     = "1"
	DEFAULT_ORG_ID_INT = 1
)

type OrgModel struct {
	Ts                int64        `json:"Ts" valid:"-"`
	OrgID             string       `json:"Id" valid:"-"`
	OrgName           string       `json:"Name" valid:"required"`
	OrgLevel          int64        `json:"OrgLevel" valid:"-"`
	SuperiorOrgID     string       `json:"SuperiorOrgId" valid:"required"`
	Comment           string       `json:"Comment" valid:"-"`
	DataOrgID         int64        `json:"DataOrgId" valid:"-"`
	Child             []*OrgModel  `json:"ChildOrgs" valid:"-"`
	OrgType           int          `json:"OrgType" valid:"-"`
	Area              int          `json:"AreaCode" valid:"-"`
	ExtData           string       `json:"ExtData" valid:"json,optional"`
	RoadType          int          `json:"RoadType" valid:"-"`
	StationCount      int          `valid:"-"`
	ErrorStationCount int          `valid:"-"`
	Users             []*UserModel `valid:"-"`
	//Status        int
}

func (org *OrgModel) IsRoot() bool {
	return org != nil && org.OrgID == DEFAULT_ORG_ID
}

type ModifyOrgModel struct {
	OrgID   string `json:"Id" valid:"required"`
	Comment string `json:"Comment" valid:"-"`
	OrgName string `json:"Name" valid:"-"`
}

func (model *ModifyOrgModel) Valid() error {
	if model.OrgName == "" && model.Comment == "" {
		return errors.New("OrgName or Comment can't be nil")
	}
	bl, err := govalidator.ValidateStruct(model)
	if bl {
		return nil
	} else {
		return err
	}
}

func (model *OrgModel) ToEntity() *entities.OrgStructure {
	entity := &entities.OrgStructure{
		Uts:           time.Now(),
		OrgName:       model.OrgName,
		OrgLevel:      model.OrgLevel,
		SuperiorOrgId: utils.FromStringToInt64(model.SuperiorOrgID),
		Comment:       model.Comment,
		DataOrgId:     model.DataOrgID,
		OrgType:       int32(model.OrgType),
		AreaCode:      int32(model.Area),
		ExtData:       model.ExtData,
		RoadType:      int32(model.RoadType),
		OrgId:         utils.FromStringToInt64(model.OrgID),
	}

	if entity.OrgId > 0 {
		entity.Ts = entity.Uts.Unix() * 1000
	} else {
		id := utils.GetIDGenerate().GetID()
		entity.OrgId = id.Int64()
		entity.Ts = id.Time()
	}
	if model.ExtData == "" {
		entity.ExtData = "{}"
	}
	//id,ts
	return entity
}

func (model *OrgModel) FromEntity(entity *entities.OrgStructure) {

	model.Comment = entity.Comment
	model.DataOrgID = entity.DataOrgId
	model.OrgID = utils.FromInt64ToString(entity.OrgId)
	model.OrgLevel = entity.OrgLevel
	model.OrgName = entity.OrgName

	model.SuperiorOrgID = utils.FromInt64ToString(entity.SuperiorOrgId)
	model.ExtData = entity.ExtData
	model.Ts = entity.Ts

}

func (model *OrgModel) Valid() error {
	bl, err := govalidator.ValidateStruct(model)
	//TODO orgsupid,orglevel
	if bl {
		return nil
	}
	return err
}

type DeleteOrgsModel struct {
	IDs []string
}

// GetID() string
// GetSuperiorID() string
// AddChild(child *TreeNode) error
// GetChild() []*TreeNode

func (model *OrgModel) GetID() string {
	return model.OrgID
}

func (model *OrgModel) GetSuperiorID() string {
	return model.SuperiorOrgID
}
func (model *OrgModel) AddChild(child interface{}) error {
	c, ok := child.(*OrgModel)
	if ok {
		model.Child = append(model.Child, c)
	}
	return nil
}
func (model *OrgModel) GetChild() interface{} {
	return model.Child
}
