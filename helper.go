// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bulk

import (
	"net"
	"sync"
	"time"
)

type interfaceInfoBase struct {
	sync.Once

	gate chan struct{}
	sync.RWMutex
	lastFetched time.Time
	toIndex     map[string]int
	toName      map[int]string
}

var ifIB interfaceInfoBase

func (ifib *interfaceInfoBase) init() {
	ifib.gate = make(chan struct{}, 1)
	ifib.toIndex = make(map[string]int)
	ifib.toName = make(map[int]string)
	ifib.fetch()
	ifib.lastFetched = time.Now()
}

func (ifib *interfaceInfoBase) fetch() {
	ift, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, ifi := range ift {
		ifib.toIndex[ifi.Name] = ifi.Index
		ifib.toName[ifi.Index] = ifi.Name
	}
}

func (ifib *interfaceInfoBase) tryAcquireSema() bool {
	select {
	case ifib.gate <- struct{}{}:
		return true
	default:
		return false
	}
}

func (ifib *interfaceInfoBase) releaseSema() {
	<-ifib.gate
}

func (ifib *interfaceInfoBase) nameToIndex(zone string) int {
	ifib.update()
	ifib.RLock()
	defer ifib.RUnlock()
	index, ok := ifib.toIndex[zone]
	if !ok {
		return 0
	}
	return index
}

func (ifib *interfaceInfoBase) indexToName(index int) string {
	ifib.update()
	ifib.RLock()
	defer ifib.RUnlock()
	name, ok := ifib.toName[index]
	if !ok {
		return ""
	}
	return name
}

func (ifib *interfaceInfoBase) update() {
	ifib.Once.Do(ifib.init)
	if !ifib.tryAcquireSema() {
		return
	}
	defer ifib.releaseSema()
	now := time.Now()
	if ifib.lastFetched.After(now.Add(-60 * time.Second)) {
		return
	}
	ifib.lastFetched = now
	ifib.Lock()
	ifib.fetch()
	ifib.Unlock()
}
