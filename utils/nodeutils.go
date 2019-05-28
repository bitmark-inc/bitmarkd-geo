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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/araujobsd/bitmarkd-geo/config"
	"github.com/araujobsd/bitmarkd-geo/geolocation"
	"github.com/flopp/go-staticmaps"
	"github.com/mmcloughlin/globe"
)

const (
	standalone = "/usr/local/bin/bitmarkd-geo-cmd"
)

func webClientIPv4() (webclient *http.Client) {
	configuration := config.LoadConfigFile()

	localAddr, err := net.ResolveIPAddr("ip", fmt.Sprintf("%s", geolocation.GetLocalIPv4(configuration["public_iface"].(string))))
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

	lat, lon, _, err = geolocation.GetLatLon(string(data))
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

func parseCsv() (rows [][]string, err error) {
	configuration := config.LoadConfigFile()

	file, err := os.Open(configuration["nodes_csv"].(string))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	rows, err = csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	return rows, err
}

func RunStandalone() (err error) {
	_, err = exec.Command(standalone).Output()
	if err != nil {
		panic(err)
	}

	return
}

func CountryTotal() (m map[string]int) {
	m = make(map[string]int)

	rows, _ := parseCsv()
	for _, val := range rows {
		valSplit := strings.Split(val[0], ",") // [0] IPAddress, [1] Country Name
		if _, ok := m[valSplit[1]]; ok {
			m[valSplit[1]]++
		} else {
			m[valSplit[1]] = 1
		}
	}

	return m
}

func NodeInCSV(ipv4 string) (lat float64, lon float64, country string, have bool) {
	rows, _ := parseCsv()
	for _, val := range rows {
		valSplit := strings.Split(val[0], ",")
		if valSplit[0] == ipv4 {
			have = true
			country = valSplit[1]
			lat, _ = strconv.ParseFloat(valSplit[2], 64)
			lon, _ = strconv.ParseFloat(valSplit[3], 64)

			return lat, lon, country, have
		}
	}

	return 0, 0, "", false
}

// WorldNodes - Creates the maps for all nodes
func WorldNodes(flatmap *sm.Context, globemap *globe.Globe, url string, m *TTLMap) (key string) {
	var lat, lon float64
	var lastKey, country string
	var have bool

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

	for _, data := range nodeInfo {
		if len(data.Listeners) > 0 {
			nodeIP := strings.FieldsFunc(data.Listeners[0], brackets)

			if strings.Contains(nodeIP[0], ".") {
				nodeIP = strings.Split(nodeIP[0], ":2136")
			}

			lat, lon, country, have = NodeInCSV(nodeIP[0])
			if have == true {
				fmt.Println("===> HAVE:", nodeIP[0])
				m.Put(nodeIP[0], country, lat, lon)
			} else {
				lat, lon, country, err = geolocation.GetLatLon(nodeIP[0])
				fmt.Println("===> DONT HAVE:", nodeIP[0])
				fmt.Println(err)
				if err == nil {
					m.Put(nodeIP[0], country, lat, lon)
				}
			}
		}
	}

	for _, v := range m.m {
		geolocation.GlobeMapAddMarker(globemap, v.Lat, v.Lon)
		geolocation.FlatMapAddMarker(flatmap, v.Lat, v.Lon)
	}

	return lastKey
}

func FindFileFlag(dir string, file []string) (flag string) {
	for _, v := range file {
		found, err := filepath.Glob(dir + v)

		if err != nil {
			fmt.Println(err)
		}

		if len(found) != 0 {
			found[0] = strings.Replace(found[0], "webserver/mysite/", "", -1)
			return found[0]
		}
	}

	return
}
