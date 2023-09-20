package api

import (
	"message/api/domain"
)

type Body struct {
	ErrorCode int              `json:"errorCode"`
	Message   string           `json:"message"`
	Data      QuerySessionResp `json:"data"`
}

type BodyV2 struct {
	ErrorCode int         `json:"errorCode"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

type QuerySessionResp struct {
	Count int64           `json:"count"`
	Array []domain.BizMsg `json:"array"`
}

type QueryStatusResp struct {
	Payload string `json:"payload"`
}

const ( // QueryStatus
	ERRORCODE_OK            = 0
	ERRORCODE__NOT_FOUND    = 1404
	ERRORCODE__INVALID_PARA = 1500
	ERRORCODE__OTHER        = 1999
)
