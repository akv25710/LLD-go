package circuit_breaker

type State string

const (
	CLOSED    State = "CLOSED"
	OPEN      State = "OPEN"
	HALF_OPEN State = "HALF-OPEN"
)
