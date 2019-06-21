// SPDX-License-Identifier: BSD-2-Clause
package utils

import (
	"net/http"
	"reflect"
	"testing"
)

func TestWebClientIPv4(t *testing.T) {
	webclient := webClientIPv4()
	if reflect.ValueOf(webclient).Kind() != reflect.Ptr {
		t.Fatal("It is not a Ptr")
	}
	if reflect.TypeOf(webclient) != reflect.TypeOf(&http.Client{}) {
		t.Fatal("It does not reflect the same http.Client type interface")
	}
}

func TestMyWanIp(t *testing.T) {
	lat, lon := MyWanIp()
	if reflect.ValueOf(lat).Kind() != reflect.Float64 ||
		reflect.ValueOf(lon).Kind() != reflect.Float64 {
		t.Fatal("It does return float64 for lat and lon")
	}
}

func TestParseNode(t *testing.T) {
	s := `[
	  {
    		"publicKey": "582e7f51251bb5a548640792d01cf6ca947436543f5f2909ab5d2ce2f3ad0f04",
    		"listeners": [
      			"127.0.0.1:1111"
    		],
    		"timestamp": "2019-05-09T10:25:59Z"
  	},
  	{
    		"publicKey": "5c19c49ad7d106189b2f70bc8a9dcec6420b79190c0883c79e8f4b0a2986b967",
    		"listeners": [
      			"127.0.0.1:1111"
    		],
    		"timestamp": "2019-05-09T10:21:49Z"
  	}
	]
	`
	nodeinfo := ParseNode([]byte(s))
	if reflect.ValueOf(nodeinfo).Kind() != reflect.Slice {
		t.Fatal("It is not a slice type")
	}
}
