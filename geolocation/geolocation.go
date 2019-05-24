/*-
 * Copyright 2019 Bitmark, Inc.
 * Copyright 2019 by Marcelo Araujo <araujo@FreeBSD.org>
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted providing that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 */

package geolocation

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
)

// GeoLocation - Get info from https://ip-api.com
type GeoLocation struct {
	As          string  `json:"as"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Isp         string  `json:"isp"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Org         string  `json:"org"`
	Query       string  `json:"query"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	Status      string  `json:"status"`
	Timezone    string  `json:"timezone"`
	Zip         string  `json:"zip"`
}

// NodeInfo - Get a list of nodes connected to each other
type NodeInfo struct {
	Listeners []string `json:"listeners"`
	PublicKey string   `json:"publicKey"`
}

var (
	ipapi      = "http://ip-api.com/json/"
	gl         GeoLocation
	CountryMap map[string]int
)

// GetLocalIPv4 - Return the local machine WAN ip
func GetLocalIPv4(iface string) (ifaceip net.IP) {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, ifc := range interfaces {
		if ifc.Name == iface {
			if addrs, err := ifc.Addrs(); err == nil {
				for _, addr := range addrs {
					switch ip := addr.(type) {
					case *net.IPNet:
						if ip.IP.DefaultMask() != nil {
							return (ip.IP)
						}
					}
				}
			} else {
				panic(err)
			}
		}
	}

	return (nil)
}

func SetCountriesNumber(countryName string) {
	if CountryMap == nil {
		CountryMap = make(map[string]int)
	}

	if _, ok := CountryMap[countryName]; ok {
		CountryMap[countryName]++
	} else {
		CountryMap[countryName] = 1
	}
}

// GetLatLon - Get lat and lon of the WAN IP
func GetLatLon(ipAddress string) (lat float64, lon float64, country string, err error) {
	url := ipapi + string(ipAddress)

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&gl)

	if gl.Status == "fail" {
		err = errors.New(gl.Status)
	}

	SetCountriesNumber(gl.Country)

	return gl.Lat, gl.Lon, gl.Country, err
}
