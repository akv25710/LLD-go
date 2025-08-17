package notification

type Status string

const (
	StatusQueued   Status = "QUEUED"
	StatusSending  Status = "SENDING"
	StatusSent     Status = "SENT"
	StatusFailed   Status = "FAILED"
	StatusRetrying Status = "RETRYING"
)
