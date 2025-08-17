package multiple

type State string

const (
	Paused    State = "Paused"
	Running         = "Running"
	Created         = "Created"
	Completed       = "Completed"
	Failed          = "Failed"
)
