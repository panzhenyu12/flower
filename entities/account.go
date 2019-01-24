package entities

import "time"

type Account struct {
	Uts           time.Time `xorm:"not null default 'now()' DATETIME"`
	Ts            int64     `xorm:"not null default 0 BIGINT"`
	UserId        int64     `xorm:"not null pk index BIGINT"`
	UserName      string    `xorm:"not null default '''::character varying' index unique VARCHAR(1024)"`
	UserPasswd    string    `xorm:"not null default '''::character varying' VARCHAR(1024)"`
	OrgId         int64     `xorm:"not null BIGINT"`
	FuncRoleId    int64     `xorm:"not null BIGINT"`
	IsValid       bool      `xorm:"not null default true BOOL"`
	RealName      string    `xorm:"default '''::character varying' VARCHAR(1024)"`
	Comment       string    `xorm:"default '''::character varying' VARCHAR"`
	Status        int       `xorm:"not null default 1 SMALLINT"`
	SecurityToken string    `xorm:"not null index unique VARCHAR(1024)"`
}

type AccountExtend struct {
	Account  `xorm:"extends"`
	FuncRole FuncRole     `xorm:"extends"`
	Org      OrgStructure `xorm:"extends"`
}
