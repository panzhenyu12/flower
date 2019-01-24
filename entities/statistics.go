package entities

import "time"

type GroupData struct {
	Type       string  `json:"type"`
	Count      int64   `json:"count"`
	PlateCount int64   `json:"plate_count"`
	AvgData    float64 `json:"avg_data,omitempty"`
	SumData    float64 `json:"sum_data,omitempty"`
}

type StatisticDay struct {
	Id       int64  `xorm:"pk" json:"-"`
	Category string `json:"-"`
	Ts       int64
	//Count     int64``
	GasStationCode string
	DateTime       time.Time    `json:"-"`
	GroupData      []*GroupData `xorm:"json" json:"-"`
}

type StatisticExtend struct {
	GroupData    `xorm:"extends"`
	StatisticDay `xorm:"extends"`
}
