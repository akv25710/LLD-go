package circuit_breaker

import (
	"errors"
	"time"
)

type Circuit struct {
	State State

	FailureThreshold int
	ResetTime        time.Duration

	Failures int
	LastRun  time.Time
}

func NewCircuit(failureThreshold int, resetTime time.Duration) *Circuit {
	return &Circuit{
		FailureThreshold: failureThreshold,
		State:            CLOSED,
		ResetTime:        resetTime,
	}
}

func (circuit *Circuit) Execute(run func() error) error {
	if circuit.State == OPEN && time.Since(circuit.LastRun) > circuit.ResetTime {
		circuit.State = HALF_OPEN
	}

	switch circuit.State {
	case OPEN:
		return errors.New("Circuit is open")
	case HALF_OPEN:
		if err := run(); err != nil {
			circuit.State = OPEN
			circuit.LastRun = time.Now()
			return err
		} else {
			circuit.State = CLOSED
		}
	case CLOSED:
		if err := run(); err != nil {
			circuit.Failures++
			if circuit.Failures >= circuit.FailureThreshold {
				circuit.State = OPEN
			}
			return err
		} else {
			circuit.Failures = 0
		}
	}

	return nil
}

func (circuit *Circuit) GetState() State {
	return circuit.State
}
