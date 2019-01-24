package model

type FuncRoleQuery struct {
	StartTimestamp int64
	EndTimestamp   int64
	Limit          int
	Offset         int
	OrderAsc       bool
	OrderBy        string
}

func (query *FuncRoleQuery) GetStartTimestamp() int64 {
	return query.StartTimestamp
}

func (query *FuncRoleQuery) GetEndTimestamp() int64 {
	return query.EndTimestamp
}

func (query *FuncRoleQuery) GetSortBy() string {
	return query.OrderBy
}

func (query *FuncRoleQuery) GetSortAsc() bool {
	return query.OrderAsc
}

func (query *FuncRoleQuery) GetLimit() int {
	return query.Limit
}

func (query *FuncRoleQuery) GetOffset() int {
	return query.Offset
}
