package api

import (
	"github.com/message/api/domain"
)

type Body struct {
	ErrorCode int              `json:"errorCode"`
	Message   string           `json:"message"`
	Data      QuerySessionResp `json:"data"`
}

type QuerySessionResp struct {
	Count int64           `json:"count"`
	Array []domain.BizMsg `json:"array"`
}
