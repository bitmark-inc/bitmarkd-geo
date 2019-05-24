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

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/araujobsd/bitmarkdgeo/geolocation"
	"github.com/araujobsd/bitmarkdgeo/utils"
)

const (
	globalTimeOut = 5
	nodesFile     = "/tmp/nodes.csv"
)

var (
	NodesCount int
	nodeUrl    = "https://node-d1.live.bitmark.com:2131/bitmarkd/peers?"
	urlCount   = "count=100"
	urlKey     = "&public_key="
	mutex      = &sync.Mutex{}
	IPlist     = make(map[string]string)
)

func dumpCSV(m *utils.TTLMap) {
	file, err := os.Create(nodesFile)
	if err != nil {
		panic("Cannot create csv file")
	}

	wr := csv.NewWriter(file)
	defer wr.Flush()

	// fmt.Sprintf("%s", geolocation.GetLocalIPv4("enp0s25")))
	mutex.Lock()
	for k, v := range m.GetAll() {
		a := strings.Fields(fmt.Sprintf("%s,%s,%f,%f", k, strings.Replace(v.Country, " ", "-", -1), v.Lat, v.Lon))
		err := wr.Write(a)
		if err != nil {
			panic("Cannot write into csv file")
		}
	}
	mutex.Unlock()
}

func main() {
	var nodeKey, lastNodeKey string
	m := utils.NewMap(globalTimeOut)

	myIPlat := 25.0478
	myIPlon := 121.5318

	flatmap := geolocation.FlatMap()
	globemap := geolocation.GlobeMap()

	if m.Len() <= 0 {
		fullUrl := nodeUrl + urlCount

		nodeKey = utils.WorldNodes(flatmap, globemap, fullUrl, m)

		for {
			if nodeKey != lastNodeKey && len(nodeKey) != 0 {
				lastNodeKey = nodeKey
				fullUrl = nodeUrl + urlCount + urlKey + lastNodeKey
				nodeKey = utils.WorldNodes(flatmap, globemap, fullUrl, m)
			} else {
				break
			}
		}

		mutex.Lock()
		geolocation.FlatMapRender(flatmap)
		geolocation.GlobeMapRender(globemap, myIPlat, myIPlon)
		mutex.Unlock()
	}

	dumpCSV(m)
}