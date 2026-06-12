package types

type QueueMetaData struct {
	Name              string `json:"name"`
	VisibilityTimeOut int32  `json:"visibilityTimeout"`
	BufferSizeLimit   int32  `json:"bufferSizeLimit"`
}
type QueueMessage struct {
	MsgId   string
	Message string
}

// CRC| RECORD LEN | ID length| ID | MSG len | MSG
