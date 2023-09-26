package apiv2

type QueryGetSessionReq struct { // to be deprecated
	//SessionId string `json:"session_id"`
	TimeAsc  bool `json:"time_asc"`
	PageId   int  `json:"page_id"`
	PageSize int  `json:"page_size"`
}

type QueryGetSessionReqForm struct {
	TimeAsc  bool `form:"time_asc"`
	PageId   int  `form:"page_id"`
	PageSize int  `form:"page_size"`
}

type Message struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
}

type QueryGetSessionResp struct {
	Obj   []interface{} `json:"obj"`
	Total int64         `json:"total"`
}
