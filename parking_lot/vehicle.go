package parking_lot

type VehicleType int

const (
	VehicleTypeCar VehicleType = iota
	VehicleTypeBike
	VehicleTypeTruck
	VehicleTypeVan
)

type Vehicle interface {
	GetId() string
	AssignTicket(parking *ParkingTicket)
	GetTicket() *ParkingTicket
	GetType() VehicleType
}

type BaseVehicle struct {
	Id     string
	Type   VehicleType
	Ticket *ParkingTicket
}

func (b *BaseVehicle) GetId() string {
	return b.Id
}

func (b *BaseVehicle) GetType() VehicleType {
	return b.Type
}

func (b *BaseVehicle) GetTicket() *ParkingTicket {
	return b.Ticket
}

func (b *BaseVehicle) AssignTicket(parking *ParkingTicket) {
	b.Ticket = parking
}

func NewVehicle(vehicleType VehicleType, id string) Vehicle {
	base := BaseVehicle{
		Id:   id,
		Type: vehicleType,
	}

	switch vehicleType {
	case VehicleTypeCar:
		return &VehicleCar{BaseVehicle: base}
	case VehicleTypeBike:
		return &VehicleBike{BaseVehicle: base}
	case VehicleTypeTruck:
		return &VehicleTruck{BaseVehicle: base}
	case VehicleTypeVan:
		return &VehicleVan{BaseVehicle: base}
	default:
		return &VehicleCar{BaseVehicle: base}
	}

}

type VehicleCar struct {
	BaseVehicle
}

type VehicleBike struct {
	BaseVehicle
}

type VehicleTruck struct {
	BaseVehicle
}
type VehicleVan struct {
	BaseVehicle
}
