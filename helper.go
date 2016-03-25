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
	indices     map[string]int
}

var ifIB interfaceInfoBase

func (ifib *interfaceInfoBase) init() {
	ifib.indices = make(map[string]int)
	ifib.fetch()
	ifib.lastFetched = time.Now()
}

func (ifib *interfaceInfoBase) fetch() {
	ift, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, ifi := range ift {
		ifib.indices[ifi.Name] = ifi.Index
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

func (ifib *interfaceInfoBase) zoneToUint32(zone string) uint32 {
	ifib.update()
	ifib.RLock()
	defer ifib.RUnlock()
	index, ok := ifib.indices[zone]
	if !ok {
		return 0
	}
	return uint32(index)
}

func (ifib *interfaceInfoBase) update() {
	ifib.Once.Do(ifib.init)
	if !ifib.tryAcquireSema() {
		return
	}
	defer ifib.releaseSema()
	now := time.Now()
	if ifib.lastFetched.After(now.Add(-5 * time.Second)) {
		return
	}
	ifib.lastFetched = now
	ifib.Lock()
	ifib.fetch()
	ifib.Unlock()
}
