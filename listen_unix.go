// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,amd64 linux,386 netbsd,amd64

package bulk

import (
	"net"
	"os"
	"syscall"
)

func listenPacket(network, address string) (*PacketConn, error) {
	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return nil, err
	}
	family := addrFamily(network, addr.IP)
	s, err := socket(family, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return nil, err
	}
	sa, err := sockaddr(family, addr.IP, addr.Port, addr.Zone)
	if err != nil {
		syscall.Close(s)
		return nil, err
	}
	if err := syscall.Bind(s, sa); err != nil {
		syscall.Close(s)
		return nil, os.NewSyscallError("bind", err)
	}
	ip, port, zone, err := getsockname(uintptr(s))
	if err != nil {
		syscall.Close(s)
		return nil, os.NewSyscallError("getsockname", err)
	}
	c := PacketConn{s: s}
	c.laddr = &net.UDPAddr{IP: ip, Port: port, Zone: zone}
	return &c, nil
}
