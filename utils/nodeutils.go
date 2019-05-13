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

package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/araujobsd/bitmarkdgeo/geolocation"
	"github.com/flopp/go-staticmaps"
	"github.com/mmcloughlin/globe"
	"github.com/schollz/progressbar"
)

func webClientIPv4() (webclient *http.Client) {
	localAddr, err := net.ResolveIPAddr("ip", fmt.Sprintf("%s", geolocation.GetLocalIPv4("enp0s25")))
	if err != nil {
		panic(err)
	}

	localTCPAddr := net.TCPAddr{
		IP: localAddr.IP,
	}

	// Set HTTP
	webcl := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: (&net.Dialer{
				LocalAddr: &localTCPAddr,
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       99 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return webcl
}

// MyWanIp - Gets the local machine WANIP and returns lat and lon.
func MyWanIp() (lat float64, lon float64) {
	webclient := webClientIPv4()
	resp, err := webclient.Get("http://myexternalip.com/raw")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	lat, lon, err = geolocation.GetLatLon(string(data))
	if err != nil {
		panic(err)
	}

	return lat, lon
}

func brackets(r rune) bool {
	return r == '[' || r == ']'
}

// ParseNode - Parses a json provide by the bitmarkd node
func ParseNode(s []byte) (nodeInfo []geolocation.NodeInfo) {
	json.Unmarshal(s, &nodeInfo)

	return nodeInfo
}

// WorldNodes - Creates the maps for all nodes
func WorldNodes(flatmap *sm.Context, globemap *globe.Globe, url string) (key string) {
	var lat, lon float64
	var lastKey string

	webclient := webClientIPv4()
	response, err := webclient.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	nodeInfo := ParseNode(data)

	if len(nodeInfo) == 0 {
		return ""
	}

	// ProgressBar
	bar := progressbar.New(10)

	for _, data := range nodeInfo {
		bar.Add(1)
		nodeIP := strings.FieldsFunc(data.Listeners[0], brackets)

		if strings.Contains(nodeIP[0], ".") {
			nodeIP = strings.Split(nodeIP[0], ":2136")
		}

		lat, lon, err = geolocation.GetLatLon(nodeIP[0])
		if err == nil {
			geolocation.GlobeMapAddMarker(globemap, lat, lon)
			geolocation.FlatMapAddMarker(flatmap, lat, lon)
		} else {
			Error.Println("Error to get information from IP:", nodeIP)
		}
		lastKey = data.PublicKey
	}

	return lastKey
}
