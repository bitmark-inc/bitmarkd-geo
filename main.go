// SPDX-License-Identifier: BSD-2-Clause
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/araujobsd/bitmarkd-geo/config"
	"github.com/araujobsd/bitmarkd-geo/utils"
)

var (
	mutex = &sync.Mutex{}
)

// Broker - It is the structure that holds clients and messages
type Broker struct {
	clients        map[chan string]bool
	newClients     chan chan string
	defunctClients chan chan string
	messages       chan string
}

// Start - It starts the routine when client is connected
func (b *Broker) Start() {
	go func() {
		for {
			select {

			case s := <-b.newClients:
				b.clients[s] = true

			case s := <-b.defunctClients:
				delete(b.clients, s)
				close(s)

			case msg := <-b.messages:
				for s := range b.clients {
					s <- msg
				}
				log.Printf("Broadcast message to %d clients", len(b.clients))
			}
		}
	}()
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Cannot flush it!", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan string)

	b.newClients <- messageChan

	notify := r.Context()
	go func() {
		<-notify.Done()
		b.defunctClients <- messageChan
		log.Println("HTTP connection closed.")
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	for {
		msg, open := <-messageChan

		if !open {
			break
		}

		fmt.Fprintf(w, "data: %s\n\n", msg)
		f.Flush()
	}

	log.Println("Finished HTTP request ", r.URL.Path)
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func main() {
	utils.InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	configuration := config.LoadConfigFile()

	b := &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}

	c := &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}

	d := &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}

	b.Start()
	c.Start()
	d.Start()
	http.Handle("/events/", b)
	http.Handle("/counter/", c)
	http.Handle("/map/", d)

	// Create flags
	go func() {
		for {
			mutex.Lock()
			countryTotal := utils.CountryTotal()
			mutex.Unlock()

			sortCountries := make([]string, 0, len(countryTotal))
			for country := range countryTotal {
				sortCountries = append(sortCountries, country)
			}
			sort.Strings(sortCountries)

			html := ""
			for _, v := range sortCountries {
				l := []string{"*" + v + "*"}
				html = html + "<div class='field_event'><span class='col col1'><center>" + strconv.Itoa(countryTotal[v]) + "</center></span>" + "<span class='col col2'><img height='30' width='40' src=" + utils.FindFileFlag("/webserver/mysite/flags/", l) + "> - " + v + "</span></div>"
			}

			b.messages <- html

			n, _ := strconv.Atoi(configuration["global_timeout"].(json.Number).String())
			time.Sleep(time.Duration(n) * time.Second)

			// Force free memory
			debug.FreeOSMemory()
		}
	}()

	// Create nodes counter
	go func() {
		for {
			mutex.Lock()
			con := utils.CountryTotal()
			mutex.Unlock()

			total := 0
			for _, v := range con {
				total += v
			}

			html := "Number of nodes: " + strconv.Itoa(total)
			c.messages <- html

			n, _ := strconv.Atoi(configuration["global_timeout"].(json.Number).String())
			time.Sleep(time.Duration(n) * time.Second)

			// Force free memory
			debug.FreeOSMemory()
		}
	}()

	// Create maps
	go func() {
		for {
			mutex.Lock()
			_ = utils.RunStandalone()
			mutex.Unlock()

			n, _ := strconv.Atoi(configuration["global_timeout"].(json.Number).String())
			time.Sleep(time.Duration(n*120) * time.Second)

			// Force free memory
			debug.FreeOSMemory()
		}
	}()

	// Add Gzip compress
	handlerNoGz := http.FileServer(http.Dir("/webserver/mysite"))
	handlerWGz := gziphandler.GzipHandler(handlerNoGz)
	http.Handle("/", handlerWGz)

	if !configuration["https"].(bool) {
		_ = http.ListenAndServe(":80", nil)
	} else {
		go http.ListenAndServe(":80", http.HandlerFunc(redirect))
		_ = http.ListenAndServeTLS(":443", "/usr/local/etc/letsencrypt/live/nodes.bitmark.com/cert.pem", "/usr/local/etc/letsencrypt/live/nodes.bitmark.com/privkey.pem", nil)
	}
}
