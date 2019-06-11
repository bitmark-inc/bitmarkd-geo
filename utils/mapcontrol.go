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
	"sync"
	"time"
)

// Nodes - This struct receives all info for each node
type Nodes struct {
	Address    string
	Country    string
	Lat        float64
	Lon        float64
	lastAccess int64
}

// TTLMap - Map with time to live
type TTLMap struct {
	m     map[string]*Nodes
	mutex sync.Mutex
}

// NewMap - Create a new map type with TTL
func NewMap(maxTTL int) (m *TTLMap) {
	m = &TTLMap{m: make(map[string]*Nodes)}
	go func() {
		for now := range time.Tick(time.Second) {
			m.mutex.Lock()
			for k, v := range m.m {
				if now.Unix()-v.lastAccess > int64(maxTTL) {
					delete(m.m, k)
				}
			}
			m.mutex.Unlock()
		}
	}()
	return
}

// Len - Return the length of a map
func (m *TTLMap) Len() int {
	return len(m.m)
}

// Put - Put a new intem into a map
func (m *TTLMap) Put(k, country string, lat float64, lon float64) {

	m.mutex.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &Nodes{Country: country, Lat: lat, Lon: lon}
		m.m[k] = it
	}
	it.lastAccess = time.Now().Unix()
	m.mutex.Unlock()
}

// Get - Get an item from a map
func (m *TTLMap) Get(k string) (v string) {
	m.mutex.Lock()
	if it, ok := m.m[k]; ok {
		v = it.Country
		it.lastAccess = time.Now().Unix()
	}
	m.mutex.Unlock()
	return

}

// GetAll - Unsafe way to get all items from a map
func (m *TTLMap) GetAll() (c map[string]*Nodes) {
	return m.m
}
