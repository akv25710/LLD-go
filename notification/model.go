package notification

import "time"

type Notification struct {
	Id          string
	RequestID   string
	UserId      int64
	Message     string
	DeviceToken string
}

type Result struct {
	Success bool
	Err     error
	Latency time.Duration
}

type Message struct {
	Id    string
	Topic string
	Data  []byte
}
