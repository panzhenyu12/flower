package model

type BaseResponse struct {
	ErrorCode int32       `json:"Code,omitempty"`
	Data      interface{} `json:"Data,omitempty"`
	ErrorMsg  string      `json:"ErrorMsg,omitempty"`
}

type ResponseData struct {
	Count int64
	Data  interface{}
}
