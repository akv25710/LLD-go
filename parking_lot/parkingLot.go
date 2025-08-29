package parking_lot

import (
	"fmt"
	"sync"
	"time"
)

type ParkingLot struct {
	mu      sync.Mutex
	Spots   []ParkingSpot
	Rates   map[ParkingType]int
	Payment *Payment
}

func (p *ParkingLot) Enter(vehicle Vehicle) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var availableSpot ParkingSpot
	for _, spot := range p.Spots {
		if spot.IsAvailableSpot() && spot.CanFit(vehicle) {
			availableSpot = spot
			break
		}
	}

	if availableSpot == nil {
		return fmt.Errorf("no spot available")
	}

	ticket := NewTicket(vehicle.GetId(), availableSpot.GetID())
	vehicle.AssignTicket(ticket)
	vehicle.SetSpot(&availableSpot)

	return availableSpot.AllotVehicle(vehicle)
}

func (p *ParkingLot) Exit(vehicle Vehicle, mode PaymentMode) (*Receipt, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	spot := vehicle.GetSpot()
	if spot == nil {
		return nil, fmt.Errorf("invalid exit vehicle %v", vehicle)
	}

	amount := p.CalculateAmount(vehicle, *spot)
	ticket := vehicle.GetTicket()
	ticket.AddExitTime()
	ticket.AddAmount(amount)
	// make payment
	if payment, err := p.Payment.Pay(amount, mode); err != nil {
		return nil, err
	} else {
		ticket.Paid = true
		err := (*spot).ReleaseVehicle(vehicle)
		if err != nil {
			return nil, err
		}
		return &Receipt{
			ID:            "REC_" + RandomAlphaNumeric(5),
			ParkingTicket: *ticket,
			Payment:       payment,
		}, nil

	}
}

func (p *ParkingLot) CalculateAmount(vehicle Vehicle, spot ParkingSpot) float64 {
	hours := time.Now().Sub(vehicle.GetTicket().InTime).Hours() + 1
	return hours * float64(p.Rates[spot.GetParkingType()])
}

var p *ParkingLot

func GetParkingLot() *ParkingLot {
	if p == nil {
		p = &ParkingLot{
			Rates: map[ParkingType]int{
				ParkingTypeBike:        10,
				ParkingTypeCompact:     20,
				ParkingTypeHandicapped: 5,
				ParkingTypeLarge:       30,
			},
			mu:      sync.Mutex{},
			Payment: NewPayment(),
		}

		// generate 100 bike spots
		for i := 1; i <= 100; i++ {
			id := "BIKE_" + RandomAlphaNumeric(6)
			p.Spots = append(p.Spots, NewParkingSpot(ParkingTypeBike, id))
		}

		for i := 1; i <= 50; i++ {
			id := "HANDI_" + RandomAlphaNumeric(6)
			p.Spots = append(p.Spots, NewParkingSpot(ParkingTypeHandicapped, id))
		}

		for i := 1; i <= 30; i++ {
			id := "LARGE_" + RandomAlphaNumeric(6)
			p.Spots = append(p.Spots, NewParkingSpot(ParkingTypeLarge, id))
		}

		for i := 1; i <= 80; i++ {
			id := "COMP_" + RandomAlphaNumeric(6)
			p.Spots = append(p.Spots, NewParkingSpot(ParkingTypeCompact, id))
		}
	}

	return p
}
