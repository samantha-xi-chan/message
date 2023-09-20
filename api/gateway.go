package api

type QueryGetSessionReq struct {
	WorkflowId string `form:"workflow_id"`
	// 套路
	SortBy   string `form:"sort_by"`
	PageId   int    `form:"page_id"`
	PageSize int    `form:"page_size"`
}

type QueryGetTaskResp struct {
	QueryGetTask []string `json:"task"`
	Total        int64    `json:"total"` // 分页之前的总数
}
