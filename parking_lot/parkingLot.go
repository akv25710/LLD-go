package parking_lot

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type ParkingLot struct {
	mu             sync.Mutex
	Spots          []ParkingSpot
	VehicleSpotMap map[string]ParkingSpot
	Rates          map[ParkingType]int
	Payment        *Payment
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

	ticket := NewTicket()
	vehicle.AssignTicket(ticket)
	if err := availableSpot.AllotVehicle(vehicle); err != nil {
		return err
	}
	p.VehicleSpotMap[vehicle.GetId()] = availableSpot

	return nil
}

func (p *ParkingLot) Exit(vehicle Vehicle, mode PaymentMode) (*Receipt, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	spot := p.VehicleSpotMap[vehicle.GetId()]
	if spot == nil {
		return nil, fmt.Errorf("invalid exit vehicle %v", vehicle)
	}

	amount := p.CalculateAmount(vehicle, spot)
	ticket := vehicle.GetTicket()
	ticket.AddExitTime()
	ticket.AddAmount(amount)
	// make payment
	if payment, err := p.Payment.Pay(amount, mode); err != nil {
		return nil, err
	} else {
		ticket.Paid = true
		spot.ReleaseVehicle(vehicle)
		return &Receipt{
			ID:            uuid.New().String(),
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
			mu:             sync.Mutex{},
			Payment:        NewPayment(),
			VehicleSpotMap: make(map[string]ParkingSpot),
		}

		// generate 100 bike spots
		for i := 1; i <= 100; i++ {
			id := "BIKE_" + uuid.New().String()
			p.Spots = append(p.Spots, NewParkingSpot(ParkingTypeBike, id))
		}

		for i := 1; i <= 50; i++ {
			id := "HANDI_" + uuid.New().String()
			p.Spots = append(p.Spots, NewParkingSpot(ParkingTypeHandicapped, id))
		}

		for i := 1; i <= 30; i++ {
			id := "LARGE_" + uuid.New().String()
			p.Spots = append(p.Spots, NewParkingSpot(ParkingTypeLarge, id))
		}

		for i := 1; i <= 80; i++ {
			id := "COMP_" + uuid.New().String()
			p.Spots = append(p.Spots, NewParkingSpot(ParkingTypeCompact, id))
		}
	}

	return p
}
