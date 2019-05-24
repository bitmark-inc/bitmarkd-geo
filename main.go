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
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/araujobsd/bitmarkdgeo/utils"
)

const (
	globalTimeOut = 5
)

var (
	nodeUrl  = "https://node-d1.live.bitmark.com:2131/bitmarkd/peers?"
	urlCount = "count=100"
	urlKey   = "&public_key="
	mutex    = &sync.Mutex{}
	IPlist   = make(map[string]string)
)

type Broker struct {
	clients        map[chan string]bool
	newClients     chan chan string
	defunctClients chan chan string
	messages       chan string
}

func (b *Broker) Start() {
	go func() {
		for {
			select {

			case s := <-b.newClients:
				b.clients[s] = true
				//log.Println("New client")

			case s := <-b.defunctClients:
				delete(b.clients, s)
				close(s)
				//log.Println("Removed client")

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

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
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

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("webserver/mysite/index.html")
	if err != nil {
		log.Fatal("Error parsing template.")

	}

	t.Execute(w, "Bitmark Inc.")
	log.Println("Finished HTTP request ", r.URL.Path)
}

func main() {
	utils.InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

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

	b.Start()
	c.Start()
	http.Handle("/events/", b)
	http.Handle("/counter/", c)

	// Create maps and flags
	go func() {
		for {
			mutex.Lock()
			_ = utils.RunStandalone()
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
				html = html + "<div class='field_event'><span class='col col1'><center>" + strconv.Itoa(countryTotal[v]) + "</center></span>" + "<span class='col col2'><img height='30' width='40' src=" + utils.FindFileFlag("webserver/mysite/flags/", l) + "> - " + v + "</span></div>"
			}

			b.messages <- fmt.Sprintf("%s", html)

			time.Sleep(time.Duration(globalTimeOut) * time.Second)
		}
	}()

	go func() {
		var html string
		var total int

		mutex.Lock()
		con := utils.CountryTotal()
		mutex.Unlock()

		for _, v := range con {
			total += v
		}

		for {
			html = "Number of nodes: " + strconv.Itoa(total)
			c.messages <- fmt.Sprintf("%s", html)

			time.Sleep(time.Duration(globalTimeOut) * time.Second)
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("webserver/mysite")))
	err := http.ListenAndServeTLS(":443", "/usr/local/etc/letsencrypt/live/nodes.bitmark.com/cert.pem", "/usr/local/etc/letsencrypt/live/nodes.bitmark.com/privkey.pem", nil)
	//err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
