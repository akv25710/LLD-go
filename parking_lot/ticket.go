package parking_lot

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type PaymentMode int

const (
	PaymentModeCash PaymentMode = iota
	PaymentModeCredit
	PaymentModeDebit
)

type PaymentStatus int

const (
	PaymentStatusPending PaymentStatus = iota
	PaymentStatusProcessing
	PaymentStatusSuccess
	PaymentStatusFailed
)

type ParkingTicket struct {
	Id      string
	Paid    bool
	InTime  time.Time
	OutTime time.Time
	Amount  float64
}

func NewTicket() *ParkingTicket {
	return &ParkingTicket{
		Id:     uuid.New().String(),
		InTime: time.Now(),
	}
}

func (t *ParkingTicket) AddExitTime() {
	t.OutTime = time.Now()
}

func (t *ParkingTicket) AddAmount(amount float64) {
	t.Amount = amount
}

type Receipt struct {
	ID            string
	ParkingTicket ParkingTicket
	Payment       *Payment
}

type Payment struct {
	Id     string
	Paid   bool
	Mode   PaymentMode
	Status PaymentStatus
}

// implements Payment interface
func (p *Payment) Pay(amount float64, mode PaymentMode) (*Payment, error) {
	// fake implementation, later you can extend logic for each mode
	if amount <= 0 {
		p.Id = uuid.New().String()
		p.Status = PaymentStatusFailed
		return p, fmt.Errorf("invalid amount")
	}

	p.Id = uuid.New().String()
	p.Mode = mode
	p.Paid = true
	p.Status = PaymentStatusSuccess
	return p, nil
}

func NewPayment(mode PaymentMode) *Payment {
	return &Payment{
		Mode: mode,
	}
}
