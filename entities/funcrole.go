package entities

import (
	"time"
)

type FuncRole struct {
	Uts          time.Time `xorm:"not null default 'now()' DATETIME"`
	Ts           int64     `xorm:"not null default 0 BIGINT"`
	FuncRoleId   int64     `xorm:"not null pk index BIGINT"`
	FuncRoleName string    `xorm:"default '''::character varying' index VARCHAR(1024)"`
	Content      string    `xorm:"not null default '''::character varying' VARCHAR(1024)"`
	Comment      string    `xorm:"default '''::character varying' VARCHAR"`
	Status       int       `xorm:"not null default 1 SMALLINT"`
	UserCount    int       `xorm:"<- INT"` // readonly
}
