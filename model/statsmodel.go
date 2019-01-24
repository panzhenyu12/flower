package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/panzhenyu12/flower/common"
	"github.com/panzhenyu12/flower/entities"
)

type OilStatsRequest struct {
	StartTimetamp int64
	EndTimestamp  int64
}

type StatsRequest struct {
	Stations       []string            `json:"Stations" valid:"required"`
	StartTimestamp int64               `json:"StartTimestamp" valid:"required"`
	EndTimestamp   int64               `json:"EndTimestamp" valid:"required"`
	TimeInterval   common.TimeInterval `json:"TimeInterval" valid:"required"`

	SumStations bool `json:"SumStations" valid:"-"`
	SumTs       bool `json:"SumTs" valid:"-"`
}

//(end int64, category, tbname string, stationcodes []string)
type StatsDalRequest struct {
	Stations        []string
	StartTimestamp  int64
	EndTimestamp    int64
	Category        string
	OriginTableName string
	StatsTableName  string
	SumStations     bool
	SumTs           bool
	//非油专用
	GroupType string //age;gender;total
	GroupByTs bool   //ts
	Origin    bool
}

// func (model *StatsRequest) GetOrgID() int64 {
// 	return utils.FromStringToInt64(model.OrgID)
// }
func (model *StatsRequest) Valid() error {
	bl, err := govalidator.ValidateStruct(model)
	if bl {
		return nil
	} else {
		return err
	}
}

//StatsResponse 统计结构
type StatsResponse struct {
	Type        string
	Count       int64
	PlateCount  int64
	Ts          int64
	StationCode string
	AvgData     float64 `json:"-"`
	SumData     float64 `json:"SumData,omitempty"`
	TimeSpan    int64   `json:"TimeSpan,omitempty"`
}

type NonOilStats struct {
	StatsResponse
	Age    string
	Gender string
}

func (model *StatsResponse) FromEntity(entity *entities.StatisticExtend) {
	model.Count = entity.Count
	model.PlateCount = entity.PlateCount
	model.Type = entity.Type
	model.Ts = entity.Ts
	model.StationCode = entity.GasStationCode
	model.AvgData = entity.AvgData
	model.SumData = entity.SumData
}

func GetStatsResponseList(datas []*entities.StatisticExtend) []*StatsResponse {
	models := make([]*StatsResponse, 0)
	for _, data := range datas {
		model := new(StatsResponse)
		model.FromEntity(data)
		models = append(models, model)
	}
	return models
}
func GetVehicleTypeList(datas []*entities.StatisticExtend) []*StatsResponse {
	models := make([]*StatsResponse, 0)

	return models
}
