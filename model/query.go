package model

type DateTimeRangeQuery interface {
	GetStartTimestamp() int64
	GetEndTimestamp() int64
}

type SortableQuery interface {
	GetSortBy() string
	GetSortAsc() bool
}

type PagingQuery interface {
	GetLimit() int
	GetOffset() int
}

type UsersQuery struct {
	BaseQuery
	OrgID string `json:"OrgId" valid:"required"`
}
type GasStationsQuery struct {
	BaseQuery
	OrgID string `json:"OrgId" valid:"required"`
}
type SensorQuery struct {
	GasstationCode string `json:"stationCode" valid:"required"`
}
type BaseQuery struct {
	StartTimestamp int64
	EndTimestamp   int64
	Desc           bool
	//PageIndex      int
	//PageSize       int
	Limit   int
	Offset  int
	OrderBy string
}

func (query BaseQuery) GetLimit() int {
	return query.Limit
}

func (query BaseQuery) GetOffset() int {
	return query.Offset
}
func (query BaseQuery) GetStartTimestamp() int64 {
	return query.StartTimestamp
}

func (query BaseQuery) GetEndTimestamp() int64 {
	return query.EndTimestamp
}

func (query BaseQuery) GetSortBy() string {
	return query.OrderBy
}

func (query BaseQuery) GetSortAsc() bool {
	return !query.Desc
}
