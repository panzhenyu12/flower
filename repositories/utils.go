package repositories

import (
	"fmt"

	"github.com/panzhenyu12/flower/model"
)

func getTimeGroupSql(timestr string) string {
	groupstr := " CAST(to_char(to_timestamp(ts/1000),'%v') as timestamp) as time_interval"
	switch timestr {
	case "hour":
		groupstr = fmt.Sprintf(groupstr, "yyyy-MM-dd HH24:00:00")
	case "day":
		groupstr = fmt.Sprintf(groupstr, "yyyy-MM-dd 00:00:00")
	default:
		groupstr = " CAST(to_char(to_timestamp(ts/1000),'yyyy-MM-dd 00:00:00') as timestamp) as time_interval"
	}
	return groupstr
}

func getWhereInSQL(ids []string, filed string) string {
	sqlwhere := ""
	if len(ids) > 0 {
		sqlwhere += fmt.Sprintf("and %v in(", filed)
		for index, id := range ids {
			sqlwhere += fmt.Sprintf("'%v'", id)
			if index+1 < len(ids) {
				sqlwhere += ","
			}
		}
		sqlwhere += ") "
	}
	return sqlwhere
}
func getWhereInSQLByParam(lens int, filed string) string {
	sqlwhere := ""

	sqlwhere += fmt.Sprintf("and %v in(", filed)
	for i := 0; i < lens; i++ {
		sqlwhere += "?"
		if i+1 < lens {
			sqlwhere += ","
		}
	}
	sqlwhere += ") "

	return sqlwhere
}

func getSqlWhereBySqlOrigin(query model.BaseQuery, param []interface{}) (string, []interface{}) {
	// param := []interface{}{id}
	sqlwhere := ""
	if query.StartTimestamp > 0 {
		sqlwhere += "and ts>=? "
		param = append(param, query.StartTimestamp)
	}
	if query.EndTimestamp > 0 {
		sqlwhere += "and ts<? "
		param = append(param, query.EndTimestamp)
	}
	return sqlwhere, param
}
