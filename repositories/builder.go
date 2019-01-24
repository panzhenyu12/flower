package repositories

import (
	"fmt"

	"github.com/go-xorm/builder"
	"github.com/panzhenyu12/flower/model"
)

type Builder struct {
	*builder.Builder
}

// alias for builder.Select
func (b *Builder) Select(cols ...string) *Builder {
	b.Builder.Select(cols...)
	return b
}

// alias for builder.From
func (b *Builder) From(subject interface{}, alias ...string) *Builder {
	b.Builder.From(subject, alias...)
	return b
}

// alias for builder.Where
func (b *Builder) Where(cond builder.Cond) *Builder {
	b.Builder.Where(cond)
	return b
}

func (b *Builder) InnerJoinById(leftTable, rightTable, col string) *Builder {
	b.InnerJoin(rightTable, fmt.Sprintf("%v.%v = %v.%v", leftTable, col, rightTable, col))
	return b
}

func (b *Builder) WhereStatusNotDeleted(tableName string) *Builder {
	b.Where(buildWhereStatusNotDeletedCond(tableName))
	return b
}

func (b *Builder) WhereTsMatch(tableName string, query model.DateTimeRangeQuery) *Builder {
	b.Builder.Where(buildWhereTsMatchCond(tableName, query))
	return b
}

func buildWhereStatusNotDeletedCond(tableName string) builder.Cond {
	return builder.In(fmt.Sprintf("%v.status", tableName), statusNotDeleted)
}

func buildWhereTsMatchCond(tableName string, query model.DateTimeRangeQuery) builder.Cond {
	cond := trueCond
	start := query.GetStartTimestamp()
	end := query.GetEndTimestamp()
	col := fmt.Sprintf("%v.ts", tableName)
	if start > 0 {
		cond = cond.And(builder.Gte{col: start})
	}
	if end > 0 {
		cond = cond.And(builder.Lte{col: end})
	}
	return cond
}

func buildWhereInCond_String(tableName, col string, values []string) builder.Cond {
	if len(values) > 0 {
		return builder.In(fmt.Sprintf("%v.%v", tableName, col), values)
	}
	return trueCond
}

func buildWhereInCond_Int(tableName, col string, values []int) builder.Cond {
	if len(values) > 0 {
		return builder.In(fmt.Sprintf("%v.%v", tableName, col), values)
	}
	return trueCond
}
func buildWhereInCond(tableName, col string, values ...interface{}) builder.Cond {
	if len(values) > 0 {
		return builder.In(fmt.Sprintf("%v.%v", tableName, col), values...)
	}
	return trueCond
}
