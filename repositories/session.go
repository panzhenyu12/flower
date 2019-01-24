package repositories

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
	"github.com/panzhenyu12/flower/common"
	"github.com/panzhenyu12/flower/model"
	"github.com/pkg/errors"
)

type Session struct {
	*xorm.Session
	engine *xorm.Engine
}

var (
	statusNotDeleted []common.TableStatus
	trueCond         = builder.And()
)

func init() {
	statusNotDeleted = []common.TableStatus{
		common.TableStatus_UnKnowStatus,
		common.TableStatus_Create,
		common.TableStatus_OK,
		common.TableStatus_Error,
	}
}

func (session *Session) Builder() *Builder {
	b := builder.Dialect(session.engine.DriverName())
	return &Builder{b}
}

func (session *Session) InnerJoin(tableName, condition string, args ...interface{}) *Session {
	session.Join("INNER", tableName, condition, args...)
	return session
}

func (session *Session) InnerJoinById(leftTable, rightTable, col string) *Session {
	return session.InnerJoin(rightTable, fmt.Sprintf("%v.%v = %v.%v", leftTable, col, rightTable, col))
}

func (session *Session) LeftJoin(tableName, condition string, args ...interface{}) *Session {
	session.Join("LEFT", tableName, condition, args...)
	return session
}

func (session *Session) LeftJoinById(leftTable, rightTable, col string) *Session {
	return session.LeftJoin(rightTable, fmt.Sprintf("%v.%v = %v.%v", leftTable, col, rightTable, col))
}
func (session *Session) LeftJoinByCode(leftTable, rightTable, leftcol string, rightcol string) *Session {
	return session.LeftJoin(rightTable, fmt.Sprintf("%v.%v = %v.%v", leftTable, leftcol, rightTable, rightcol))
}

//builder.Cond
func (session *Session) WhereCond(cond builder.Cond) *Session {
	session.Where(cond)
	return session
}
func (session *Session) WhereBetween(col string, lessVal, moreVal interface{}) *Session {
	session.Where(builder.Between{col, lessVal, moreVal})
	return session
}

func (session *Session) WhereLike(col, val string) *Session {
	session.Where(builder.Like{col, val})
	return session
}

func (session *Session) WhereEq(col string, val interface{}) *Session {
	session.Where(builder.Eq{col: val})
	return session
}

func (session *Session) WhereStatusNotDeleted(tableName string) *Session {
	session.Where(buildWhereStatusNotDeletedCond(tableName))
	return session
}

func WhereStatusNotDeletedSQL(tableName string) string {
	arr := make([]string, len(statusNotDeleted))
	for i := 0; i < len(statusNotDeleted); i++ {
		arr[i] = fmt.Sprintf("%d", statusNotDeleted[i])
	}
	return fmt.Sprintf("%v.status IN (%v)", tableName, strings.Join(arr, ","))
}

func (session *Session) WhereTsMatch(tableName string, query model.DateTimeRangeQuery) *Session {
	session.Where(buildWhereTsMatchCond(tableName, query))
	return session
}
func (session *Session) OrderBy(tableName string, orderby string, isAsc bool) *Session {
	if orderby == "" {
		return session
	}
	col := fmt.Sprintf("%v.%v", tableName, orderby)
	if isAsc {
		session.Asc(col)
	} else {
		session.Desc(col)
	}
	return session
}

func (session *Session) Sort(tableName string, query model.SortableQuery) *Session {
	sortBy := query.GetSortBy()
	if sortBy == "" {
		return session
	}
	col := fmt.Sprintf("%v.%v", tableName, sortBy)
	if query.GetSortAsc() {
		session.Asc(col)
	} else {
		session.Desc(col)
	}
	return session
}

func (session *Session) Page(query model.PagingQuery) *Session {
	if limit := query.GetLimit(); limit > 0 {
		session.Limit(limit, query.GetOffset())
	}
	return session
}

// NOTE: This method DO NOT escape values
func (this *Session) BatchUpdateField(tableName, field, idField string, values map[string]interface{}) error {
	valueStrs := make([]string, 0)
	for id, value := range values {
		var buffer bytes.Buffer
		if _, err := buffer.WriteString("('"); err != nil {
			return errors.WithStack(err)
		}
		if _, err := buffer.WriteString(id); err != nil {
			return errors.WithStack(err)
		}
		if _, err := buffer.WriteString("',"); err != nil {
			return errors.WithStack(err)
		}
		if str, ok := value.(string); ok {
			if _, err := buffer.WriteString("'"); err != nil {
				return errors.WithStack(err)
			}
			// if _, err := buffer.WriteString(escape(str)); err != nil {
			if _, err := buffer.WriteString(str); err != nil {
				return errors.WithStack(err)
			}
			if _, err := buffer.WriteString("'"); err != nil {
				return errors.WithStack(err)
			}
		} else {
			if _, err := buffer.WriteString(fmt.Sprintf("%v", value)); err != nil {
				return errors.WithStack(err)
			}
		}
		if _, err := buffer.WriteString(fmt.Sprintf(")")); err != nil {
			return errors.WithStack(err)
		}
		valueStrs = append(valueStrs, buffer.String())
	}
	sql := fmt.Sprintf("UPDATE %v SET %v = v.val FROM (VALUES %v ) AS v (id, val) WHERE %v.%v = v.id", tableName, field, strings.Join(valueStrs, ","), tableName, idField)
	_, err := this.Exec(sql)
	return errors.WithStack(err)
}

// NOTE: This method DO NOT escape values
func (this *Session) BatchUpdateFields(tableName string, fields []string, values [][]interface{}) error {
	if len(fields) < 2 {
		return fmt.Errorf("Invalid fields")
	}
	if len(values) == 0 {
		return fmt.Errorf("Invalid values")
	}
	idField := fields[0]
	valueStrs := make([]string, len(values))
	for i, value := range values {
		if len(value) != len(fields) {
			return fmt.Errorf("Invalid value")
		}
		strs := make([]string, len(value))
		for j, val := range value {
			if str, ok := val.(string); ok {
				strs[j] = fmt.Sprintf("'%v'", str)
			} else {
				strs[j] = fmt.Sprintf("%v", val)
			}
		}
		valueStrs[i] = fmt.Sprintf("(%v)", strings.Join(strs, ","))
	}
	setStrs := make([]string, 0)
	for _, field := range fields[1:] {
		setStrs = append(setStrs, fmt.Sprintf("%v = v.%v", field, field))
	}
	sql := fmt.Sprintf("UPDATE %v SET %v FROM (VALUES %v ) AS v (%v) WHERE %v.%v = v.%v",
		tableName,
		strings.Join(setStrs, ","),
		strings.Join(valueStrs, ","),
		strings.Join(fields, ","),
		tableName, idField, idField)
	_, err := this.Exec(sql)
	return errors.WithStack(err)
}

func (this *Session) EstimateCount(query string) (int, error) {
	sql := fmt.Sprintf("SELECT count_estimate('%v')", strings.Replace(query, "'", "''", -1))
	var rows []int
	err := this.SQL(sql).Find(&rows)
	if err != nil {
		return 0, errors.WithStack(err)
	} else if len(rows) != 1 {
		return 0, fmt.Errorf("Failed to get rows %v", rows)
	}
	return rows[0], nil
}

func buildOrderBy(tableName string, query model.SortableQuery) string {
	sortBy := query.GetSortBy()
	if sortBy == "" {
		return ""
	}
	var order string
	if query.GetSortAsc() {
		order = "ASC"
	} else {
		order = "DESC"
	}
	return fmt.Sprintf("%v.%v %v", tableName, sortBy, order)
}
