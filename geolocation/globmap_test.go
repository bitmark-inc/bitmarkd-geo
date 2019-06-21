// SPDX-License-Identifier: BSD-2-Clause
package geolocation

import (
	"reflect"
	"testing"

	"github.com/flopp/go-staticmaps"
	"github.com/mmcloughlin/globe"
)

func TestFlatMap(t *testing.T) {
	flatmap := FlatMap()
	if reflect.TypeOf(flatmap) != reflect.TypeOf(&sm.Context{}) {
		t.Fatal("It does not reflect the same sm.Context type interface")
	}

	if reflect.ValueOf(flatmap).Elem().FieldByName("width").Int() != 1920 {
		t.Error("Width different than 1920")
	}
	if reflect.ValueOf(flatmap).Elem().FieldByName("height").Int() != 1080 {
		t.Error("Height different than 1080")
	}
	if reflect.ValueOf(flatmap).Elem().FieldByName("hasCenter").Bool() != true {
		t.Error("Center is not true")
	}
}

func TestGlobeMap(t *testing.T) {
	globmap := GlobeMap()
	if reflect.TypeOf(globmap) != reflect.TypeOf(&globe.Globe{}) {
		t.Fatal("It does not reflect the same globe.Globe type interface")
	}
}
