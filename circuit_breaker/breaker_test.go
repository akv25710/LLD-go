package circuit_breaker

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_OpenOnFailures(t *testing.T) {
	c := NewCircuit(3, 1*time.Second)

	// Force failures to trip breaker
	failFn := func() error { return errors.New("failure") }
	for i := 0; i < 3; i++ {
		_ = c.Execute(failFn)
	}

	if c.GetState() != OPEN {
		t.Errorf("expected state=OPEN, got %v", c.GetState())
	}
}

func TestCircuitBreaker_ResetsAfterTimeout(t *testing.T) {
	c := NewCircuit(1, 500*time.Millisecond)

	failFn := func() error { return errors.New("failure") }
	_ = c.Execute(failFn) // open circuit

	if c.GetState() != OPEN {
		t.Fatalf("expected state=OPEN, got %v", c.GetState())
	}

	time.Sleep(600 * time.Millisecond) // wait for reset window

	successFn := func() error { return nil }
	_ = c.Execute(successFn) // should go half-open -> closed

	if c.GetState() != CLOSED {
		t.Errorf("expected state=CLOSED after successful retry, got %v", c.GetState())
	}
}

func TestCircuitBreaker_HalfOpenFailureGoesOpen(t *testing.T) {
	c := NewCircuit(1, 200*time.Millisecond)

	failFn := func() error { return errors.New("failure") }
	_ = c.Execute(failFn) // open circuit

	time.Sleep(250 * time.Millisecond) // let it move to HALF-OPEN

	// Fail again in HALF-OPEN
	_ = c.Execute(failFn)

	if c.GetState() != OPEN {
		t.Errorf("expected state=OPEN after half-open failure, got %v", c.GetState())
	}
}

func TestCircuitBreaker_BlocksWhenOpen(t *testing.T) {
	c := NewCircuit(1, 2*time.Second)

	failFn := func() error { return errors.New("failure") }
	_ = c.Execute(failFn) // should trip to OPEN

	if c.GetState() != OPEN {
		t.Fatalf("expected state=OPEN, got %v", c.GetState())
	}

	successFn := func() error { return nil }

	// First call after open → still executes normally
	_ = c.Execute(successFn)

	// Second call after open → should now block
	err := c.Execute(successFn)
	if err == nil || err.Error() != "Circuit is open" {
		t.Errorf("expected circuit to block, got %v", err)
	}
}
