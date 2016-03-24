// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64

package bulk

import (
	"net"
	"syscall"
	"unsafe"
)

const (
	sysRECVMMSG = 299 // 337 for arm, 0x40000219 for 386
	sysSENDMMSG = 307 // 345 for arm, 0x4000021a for 386
)

func msgSockaddr(ip net.IP, port int) (*byte, uint32) {
	if ip.To4() != nil {
		sa := sysSockaddrInet{Family: syscall.AF_INET}
		p := (*[2]byte)(unsafe.Pointer(&sa.Port))
		p[0], p[1] = byte(port>>8), byte(port)
		copy(sa.Addr[:], ip.To4())
		return (*byte)(unsafe.Pointer(&sa)), sysSizeofSockaddrInet
	}
	if ip.To16() != nil && ip.To4() == nil {
		sa := sysSockaddrInet6{Family: syscall.AF_INET6}
		p := (*[2]byte)(unsafe.Pointer(&sa.Port))
		p[0], p[1] = byte(port>>8), byte(port)
		copy(sa.Addr[:], ip)
		return (*byte)(unsafe.Pointer(&sa)), sysSizeofSockaddrInet6
	}
	return nil, 0
}
