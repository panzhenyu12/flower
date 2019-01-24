package entities

import (
	"time"
)

type OrgStructure struct {
	Uts           time.Time
	Ts            int64
	OrgId         int64 `xorm:"pk"`
	OrgName       string
	OrgLevel      int64
	SuperiorOrgId int64
	Comment       string
	Status        int32
	DataOrgId     int64
	OrgType       int32
	AreaCode      int32
	ExtData       string
	RoadType      int32
}
type OrgStructureExtends struct {
	OrgStructure `xorm:"extends"`
	IsOrg        bool
	Code         string
}
