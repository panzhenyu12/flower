package repositories

import (
	"fmt"

	"github.com/golang/glog"
	//"flower/common"
	"flower/entities"
	"flower/model"
	"flower/utils"
)

type StatisticDal struct {
	*DB
}

func GetStatisticDal() *StatisticDal {
	return &StatisticDal{
		DB: GetDeepDB(),
	}
}

func isNonOil(req model.StatsDalRequest) bool {

	return false
}

func (dal *StatisticDal) GetTodayStatsTemplet(req model.StatsDalRequest, fcday, fchour func(req model.StatsDalRequest) ([]*entities.StatisticExtend, error)) ([]*entities.StatisticExtend, error) {
	starttshour := utils.GetHourTime()
	starttsday := utils.GetDayTime()
	ts := utils.TimeToTimestamp(starttshour)
	req.StartTimestamp = utils.TimeToTimestamp(starttsday)
	req.StatsTableName = ""
	data, err := fcday(req)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	if req.EndTimestamp > ts {
		req.StartTimestamp = ts
		if hourdata, err := fchour(req); err == nil {
			maphdata := make(map[string]*entities.StatisticExtend, 0)
			for _, d := range hourdata {
				d.Ts = utils.TimeToTimestamp(starttsday)
				maphdata[d.GasStationCode+d.Type] = d
			}
			for _, d := range data {
				d.Ts = utils.TimeToTimestamp(starttsday)
				key := d.GasStationCode + d.Type
				if m, ok := maphdata[key]; ok {
					d.Count += m.Count
					d.PlateCount += m.PlateCount
					d.SumData += m.SumData
					if d.SumData > 0 && d.Count > 0 {
						d.AvgData = d.SumData / float64(d.Count)
					}
					delete(maphdata, key)
				}
			}
			for _, v := range maphdata {
				data = append(data, v)
			}
		} else {
			glog.Error(err)
		}

	}
	return data, nil
}
func (dal *StatisticDal) GetStatsData(req model.StatsDalRequest) ([]*entities.StatisticExtend, error) {
	if isNonOil(req) {
		//req.GroupByTs = !req.SumTs
		return dal.GetNonOilStatsData(req)
	}
	return dal.GetOilStatsData(req)
}

//GetOilStatsData 获取历史统计信息
func (dal *StatisticDal) GetOilStatsData(req model.StatsDalRequest) ([]*entities.StatisticExtend, error) {
	sql := `SELECT * FROM %v, jsonb_populate_recordset(null::stats_group_data,group_data)
			where category=? and ts>=? and ts<? %v`
	param := []interface{}{req.Category, req.StartTimestamp, req.EndTimestamp}
	stationcodes := req.Stations
	sqlwhere := ""
	if stationcodes != nil && len(stationcodes) > 0 {
		sqlwhere = getWhereInSQLByParam(len(stationcodes), "gas_station_code")
		for _, code := range stationcodes {
			param = append(param, code)
		}
	}
	data := make([]*entities.StatisticExtend, 0)
	sql = fmt.Sprintf(sql, req.StatsTableName, sqlwhere)
	sqlts := ""
	if req.SumTs == false {
		sqlts = ",ts"
	}
	if req.SumStations {
		sql = fmt.Sprintf(`select type ,sum(count)as count ,sum(plate_count) as plate_count,avg(avg_data) as avg_data ,sum(sum_data) as sum_data %v from(
		%v ) as foo group by type %v`, sqlts, sql, sqlts)
	}
	err := dal.Engine.Sql(sql, param...).Find(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//GetNonOilStatsData 非油统计，根据年龄，性别，总数进行sum.
func (dal *StatisticDal) GetNonOilStatsData(req model.StatsDalRequest) ([]*entities.StatisticExtend, error) {
	sql := `select f.*,f.t as type from (
		select sum(count) as count,%v as t %v %v from (
		select  (regexp_split_to_array(type,'_'))[1]  as age,(regexp_split_to_array(type,'_'))[2] as gender,'total' as total,* 
		from %v,jsonb_populate_recordset(null::stats_group_data,group_data)
		 where category=? and ts>=? and ts<? %v
		) as foo group by t %v %v
		) as f`
	param := []interface{}{req.Category, req.StartTimestamp, req.EndTimestamp}
	stationcodes := req.Stations
	sqlwhere := ""
	if stationcodes != nil && len(stationcodes) > 0 {
		sqlwhere = getWhereInSQLByParam(len(stationcodes), "gas_station_code")
		for _, code := range stationcodes {
			param = append(param, code)
		}
	}
	sqlts := ""
	if req.SumTs == false {
		sqlts = ",ts"
	}
	sqlstation := ""
	if !req.SumStations {
		sqlstation = ",gas_station_code"
	}
	data := make([]*entities.StatisticExtend, 0)
	sql = fmt.Sprintf(sql, req.GroupType, sqlstation, sqlts, req.StatsTableName, sqlwhere, sqlstation, sqlts)

	err := dal.Engine.Sql(sql, param...).Find(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//SumStatsCount sum统计表数据，根据type
func (dal *StatisticDal) SumStatsCount(req model.StatsDalRequest) ([]*entities.StatisticExtend, error) {
	sql := `select category %v,type,sum(count) as count,sum(plate_count) as plate_count,avg(avg_data) as avg_data ,sum(sum_data) as sum_data from %v,jsonb_populate_recordset(null::stats_group_data,group_data) as t
	where category=? and ts>=? and ts<? %v group by category %v,type`
	param := []interface{}{req.Category, req.StartTimestamp, req.EndTimestamp}
	sqlwhere := ""
	stationcodes := req.Stations
	if stationcodes != nil && len(stationcodes) > 0 {
		sqlwhere = getWhereInSQLByParam(len(stationcodes), "gas_station_code")
		for _, code := range stationcodes {
			param = append(param, code)
		}
	}
	data := make([]*entities.StatisticExtend, 0)
	sqlstation := ""
	if !req.SumStations {
		sqlstation = ",gas_station_code"
	}
	sql = fmt.Sprintf(sql, sqlstation, req.StatsTableName, sqlwhere, sqlstation)
	err := dal.Engine.Sql(sql, param...).Find(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//GetTodayStatData 所有类型今日统计（不包含目前一小时）
func (dal *StatisticDal) GetTodayStatData(req model.StatsDalRequest, fc func(req model.StatsDalRequest) ([]*entities.StatisticExtend, error)) ([]*entities.StatisticExtend, error) {
	return dal.GetTodayStatsTemplet(req, dal.SumStatsCount, fc)
}
