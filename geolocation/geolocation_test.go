package geolocation

import (
	"reflect"
	"testing"
)

func TestGetLocalIPv4(t *testing.T) {
	ifaceip := GetLocalIPv4("virbr0")

	if ifaceip == nil {
		t.Error("Could not find virbr0 interface, well, nothing critical here")
	}
}

func TestGetLatLon(t *testing.T) {
	ipAddress := "8.8.8.8"
	lat, lon, err := GetLatLon(ipAddress)

	if err != nil {
		t.Error("Probably you have no Internet")
	}

	if reflect.ValueOf(lat).Kind() != reflect.Float64 {
		t.Fatal("The value of lat must be float64")
	}

	if reflect.ValueOf(lon).Kind() != reflect.Float64 {
		t.Fatal("The value of lon must be float64")
	}
}
