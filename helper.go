// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bulk

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type ipv6ZoneCache struct {
	sync.RWMutex
	lastFetched time.Time
	toIndex     map[string]int
	toName      map[int]string
}

var zoneCache = ipv6ZoneCache{
	toIndex: make(map[string]int),
	toName:  make(map[int]string),
}

func (zc *ipv6ZoneCache) nameToIndex(name string) int {
	if name == "" {
		return 0
	}
	zc.update()
	zc.RLock()
	defer zc.RUnlock()
	index, ok := zc.toIndex[name]
	if !ok {
		index, _ = strconv.Atoi(name)
	}
	return index
}

func (zc *ipv6ZoneCache) indexToName(index int) string {
	if index == 0 {
		return ""
	}
	zc.update()
	zc.RLock()
	defer zc.RUnlock()
	name, ok := zc.toName[index]
	if !ok {
		name = fmt.Sprintf("%d", index)
	}
	return name
}

func (zc *ipv6ZoneCache) update() {
	zc.Lock()
	defer zc.Unlock()
	now := time.Now()
	if zc.lastFetched.After(now.Add(-60 * time.Second)) {
		return
	}
	zc.lastFetched = now
	ift, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, ifi := range ift {
		delete(zc.toIndex, ifi.Name)
		zc.toIndex[ifi.Name] = ifi.Index
		delete(zc.toName, ifi.Index)
		zc.toName[ifi.Index] = ifi.Name
	}
}
