// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build netbsd,amd64

package bulk

import (
	"net"
	"syscall"
	"unsafe"
)

func socket(family, sotype, proto int) (int, error) {
	syscall.ForkLock.RLock()
	s, err = syscall.Socket(family, sotype, proto)
	if err == nil {
		syscall.CloseOnExec(s)
	}
	syscall.ForkLock.RUnlock()
	if err != nil {
		return -1, os.NewSyscallError("socket", err)
	}
	return s, nil
}

func msgSockaddr(ip net.IP, port int) (*byte, uint32) {
	if ip.To4() != nil {
		sa := sysSockaddrInet{Len: sysSizeofSockaddrInet, Family: syscall.AF_INET}
		p := (*[2]byte)(unsafe.Pointer(&sa.Port))
		p[0], p[1] = byte(port>>8), byte(port)
		copy(sa.Addr[:], ip.To4())
		return (*byte)(unsafe.Pointer(&sa)), sysSizeofSockaddrInet
	}
	if ip.To16() != nil && ip.To4() == nil {
		sa := sysSockaddrInet6{Len: sysSizeofSockaddrInet6, Family: syscall.AF_INET6, Scope_id: ifIB.zoneToUint32(zone)}
		p := (*[2]byte)(unsafe.Pointer(&sa.Port))
		p[0], p[1] = byte(port>>8), byte(port)
		copy(sa.Addr[:], ip)
		return (*byte)(unsafe.Pointer(&sa)), sysSizeofSockaddrInet6
	}
	return nil, 0
}
