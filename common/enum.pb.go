// Code generated by protoc-gen-go. DO NOT EDIT.
// source: enum.proto

package common

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

//表状态枚举
type TableStatus int32

const (
	//未知
	TableStatus_UnKnowStatus TableStatus = 0
	//创建
	TableStatus_Create TableStatus = 1
	//正常
	TableStatus_OK TableStatus = 2
	//错误
	TableStatus_Error TableStatus = 3
	//删除
	TableStatus_Delete TableStatus = 4
)

var TableStatus_name = map[int32]string{
	0: "UnKnowStatus",
	1: "Create",
	2: "OK",
	3: "Error",
	4: "Delete",
}

var TableStatus_value = map[string]int32{
	"UnKnowStatus": 0,
	"Create":       1,
	"OK":           2,
	"Error":        3,
	"Delete":       4,
}

func (x TableStatus) String() string {
	return proto.EnumName(TableStatus_name, int32(x))
}

func (TableStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_13a9f1b5947140c8, []int{0}
}

type ErrorCode int32

const (
	ErrorCode_UnKnowErrorCode ErrorCode = 0
	//未找到
	ErrorCode_NotFound ErrorCode = 1
	//数据库错误
	ErrorCode_DBError ErrorCode = 2
	//不能操作
	ErrorCode_NoOperate ErrorCode = 3
	//以经在操作了
	ErrorCode_Busy ErrorCode = 4
)

var ErrorCode_name = map[int32]string{
	0: "UnKnowErrorCode",
	1: "NotFound",
	2: "DBError",
	3: "NoOperate",
	4: "Busy",
}

var ErrorCode_value = map[string]int32{
	"UnKnowErrorCode": 0,
	"NotFound":        1,
	"DBError":         2,
	"NoOperate":       3,
	"Busy":            4,
}

func (x ErrorCode) String() string {
	return proto.EnumName(ErrorCode_name, int32(x))
}

func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_13a9f1b5947140c8, []int{1}
}

type TimeInterval int32

const (
	TimeInterval_UnKnowTimeInterval TimeInterval = 0
	TimeInterval_Hour               TimeInterval = 1
	TimeInterval_Day                TimeInterval = 2
)

var TimeInterval_name = map[int32]string{
	0: "UnKnowTimeInterval",
	1: "Hour",
	2: "Day",
}

var TimeInterval_value = map[string]int32{
	"UnKnowTimeInterval": 0,
	"Hour":               1,
	"Day":                2,
}

func (x TimeInterval) String() string {
	return proto.EnumName(TimeInterval_name, int32(x))
}

func (TimeInterval) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_13a9f1b5947140c8, []int{2}
}

func init() {
	proto.RegisterEnum("common.TableStatus", TableStatus_name, TableStatus_value)
	proto.RegisterEnum("common.ErrorCode", ErrorCode_name, ErrorCode_value)
	proto.RegisterEnum("common.TimeInterval", TimeInterval_name, TimeInterval_value)
}

func init() { proto.RegisterFile("enum.proto", fileDescriptor_13a9f1b5947140c8) }

var fileDescriptor_13a9f1b5947140c8 = []byte{
	// 212 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0xcd, 0x4a, 0xc3, 0x50,
	0x10, 0x85, 0xf3, 0x67, 0xda, 0x9c, 0x46, 0x1c, 0x46, 0xf0, 0x21, 0xb2, 0x70, 0xe3, 0xca, 0x6d,
	0x1b, 0x45, 0x2d, 0xb4, 0x0b, 0xe3, 0x03, 0xdc, 0xda, 0x59, 0x08, 0xc9, 0x9d, 0x32, 0xbd, 0x57,
	0xe9, 0xdb, 0x4b, 0x12, 0x10, 0x97, 0xf3, 0x0d, 0xe7, 0x3b, 0x1c, 0x40, 0x7c, 0x1c, 0xee, 0x4f,
	0xa6, 0x41, 0xb9, 0xfc, 0xd4, 0x61, 0x50, 0xdf, 0xbc, 0x61, 0xd5, 0xb9, 0x43, 0x2f, 0xef, 0xc1,
	0x85, 0x78, 0x66, 0x42, 0xfd, 0xe1, 0xb7, 0x5e, 0x7f, 0xe6, 0x9b, 0x12, 0x06, 0xca, 0x8d, 0x89,
	0x0b, 0x42, 0x29, 0x97, 0xc8, 0xf6, 0x5b, 0xca, 0xb8, 0xc2, 0xd5, 0x93, 0x99, 0x1a, 0xe5, 0xe3,
	0xbb, 0x95, 0x5e, 0x82, 0x50, 0xd1, 0x74, 0xa8, 0x26, 0xbc, 0xd1, 0xa3, 0xf0, 0x2d, 0x6e, 0x66,
	0xd3, 0x1f, 0xa2, 0x84, 0x6b, 0x2c, 0x77, 0x1a, 0x9e, 0x35, 0xfa, 0x23, 0xa5, 0xbc, 0xc2, 0xa2,
	0x5d, 0xcf, 0xa2, 0x8c, 0xaf, 0x51, 0xed, 0x74, 0x7f, 0x12, 0x1b, 0xab, 0x72, 0x5e, 0xa2, 0x58,
	0xc7, 0xf3, 0x85, 0x8a, 0xe6, 0x11, 0x75, 0xf7, 0x35, 0xc8, 0xab, 0x0f, 0x62, 0xdf, 0xae, 0xe7,
	0x3b, 0xf0, 0x2c, 0xfe, 0x4f, 0x29, 0x19, 0x13, 0x2f, 0x1a, 0x8d, 0x52, 0x5e, 0x20, 0x6f, 0xdd,
	0x85, 0xb2, 0x43, 0x39, 0x6d, 0x7d, 0xf8, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xf0, 0xcd, 0xbd, 0x71,
	0xf9, 0x00, 0x00, 0x00,
}
