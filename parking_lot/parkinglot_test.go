package parking_lot

import "testing"

func TestParkingLot_Enter_Exit(t *testing.T) {
	// write code to test
	parkingLot := GetParkingLot()
	vehicle := NewVehicle(VehicleTypeCar, "BR9F0372")
	if err := parkingLot.Enter(vehicle); err != nil {
		t.Error(err)
	} else {
		t.Logf("vehcile entered %v", vehicle)
	}

	if rec, err := parkingLot.Exit(vehicle, PaymentModeCash); err != nil {
		t.Error(err)
	} else {
		t.Logf("vehcile leaved %v", rec)
	}
}
