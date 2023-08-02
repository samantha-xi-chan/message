package domain

const (
	SESSION_START = 1001
	SESSION_END   = 1002
	SESSION_ING   = 1003
)

type UpdateSessionStatus struct { // 任务启动、关闭
	SessionID string      `json:"session_id"`
	Timestamp int64       `json:"timestamp"`
	EvtType   int32       `json:"event_type"`
	Payload   interface{} `json:"payload"`
}

type FeedSessionStream struct { // 任务内部 新的报文
	SessionID string      `json:"session_id"`
	Timestamp int64       `json:"timestamp"`
	EvtType   int32       `json:"event_type"`
	Payload   interface{} `json:"payload"`
}

type FeedSessionStreamDAO struct {
	SessionID string      `json:"session_id"`
	Timestamp int64       `json:"timestamp"`
	Payload   interface{} `json:"payload"`
	Deleted   bool        `json:"deleted"`
}
