package parking_lot

type ParkingType int

const (
	ParkingTypeHandicapped ParkingType = iota
	ParkingTypeCompact
	ParkingTypeLarge
	ParkingTypeBike
)

type ParkingSpot interface {
	GetID() string
	IsAvailableSpot() bool
	GetParkingType() ParkingType
	CanFit(v Vehicle) bool
	AllotVehicle(v Vehicle) error
	ReleaseVehicle(v Vehicle) error
}

type BaseParkingSpot struct {
	ID          string
	IsAvailable bool
	ParkingType ParkingType
	vehicle     Vehicle
}

func (p *BaseParkingSpot) GetID() string {
	return p.ID
}

func (p *BaseParkingSpot) GetParkingType() ParkingType {
	return p.ParkingType
}

func (p *BaseParkingSpot) IsAvailableSpot() bool {
	return p.IsAvailable
}

func (p *BaseParkingSpot) AllotVehicle(v Vehicle) error {
	p.vehicle = v
	p.IsAvailable = false
	return nil
}

func (p *BaseParkingSpot) ReleaseVehicle(v Vehicle) error {
	p.vehicle = nil
	p.IsAvailable = true
	return nil
}

func NewParkingSpot(pt ParkingType, id string) ParkingSpot {
	base := BaseParkingSpot{
		IsAvailable: true,
		ParkingType: pt,
		ID:          id,
	}

	switch pt {
	case ParkingTypeHandicapped:
		return &ParkingSpotHandi{BaseParkingSpot: base}
	case ParkingTypeCompact:
		return &ParkingSpotCompact{base}
	case ParkingTypeLarge:
		return &ParkingSpotLarge{base}
	case ParkingTypeBike:
		return &ParkingSpotBike{base}
	default:
		return &ParkingSpotLarge{base}
	}

}

type ParkingSpotHandi struct {
	BaseParkingSpot
}

func (p *ParkingSpotHandi) CanFit(v Vehicle) bool {
	return v.GetType() == VehicleTypeBike || v.GetType() == VehicleTypeVan || v.GetType() == VehicleTypeCar
}

type ParkingSpotCompact struct {
	BaseParkingSpot
}

func (p *ParkingSpotCompact) CanFit(v Vehicle) bool {
	return v.GetType() == VehicleTypeBike || v.GetType() == VehicleTypeVan || v.GetType() == VehicleTypeCar
}

type ParkingSpotLarge struct {
	BaseParkingSpot
}

func (p *ParkingSpotLarge) CanFit(v Vehicle) bool {
	return v.GetType() == VehicleTypeTruck || v.GetType() == VehicleTypeBike || v.GetType() == VehicleTypeVan || v.GetType() == VehicleTypeCar
}

type ParkingSpotBike struct {
	BaseParkingSpot
}

func (p *ParkingSpotBike) CanFit(v Vehicle) bool {
	return v.GetType() == VehicleTypeBike
}
