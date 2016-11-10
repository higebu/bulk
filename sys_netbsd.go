// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build netbsd,amd64 netbsd,386

package bulk

import (
	"net"
	"os"
	"syscall"
	"unsafe"
)

func socket(family, sotype, proto int) (uintptr, error) {
	syscall.ForkLock.RLock()
	s, err := syscall.Socket(family, sotype, proto)
	if err == nil {
		syscall.CloseOnExec(s)
	}
	syscall.ForkLock.RUnlock()
	if err != nil {
		return ^uintptr(0), os.NewSyscallError("socket", err)
	}
	return uintptr(s), nil
}

func soclose(s uintptr) error { return syscall.Close(int(s)) }

func msgSockaddr(ip net.IP, port int, zone string) (*byte, uint32) {
	if ip.To4() != nil {
		sa := sockaddrInet{Len: sizeofSockaddrInet, Family: syscall.AF_INET}
		p := (*[2]byte)(unsafe.Pointer(&sa.Port))
		p[0], p[1] = byte(port>>8), byte(port)
		copy(sa.Addr[:], ip.To4())
		return (*byte)(unsafe.Pointer(&sa)), sizeofSockaddrInet
	}
	if ip.To16() != nil && ip.To4() == nil {
		sa := sockaddrInet6{Len: sizeofSockaddrInet6, Family: syscall.AF_INET6, Scope_id: uint32(zoneCache.index(zone))}
		p := (*[2]byte)(unsafe.Pointer(&sa.Port))
		p[0], p[1] = byte(port>>8), byte(port)
		copy(sa.Addr[:], ip)
		return (*byte)(unsafe.Pointer(&sa)), sizeofSockaddrInet6
	}
	return nil, 0
}
