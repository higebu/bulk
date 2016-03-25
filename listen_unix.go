// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux netbsd

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
	var family int
	if addr.IP.To4() != nil {
		family = syscall.AF_INET
	}
	if addr.IP.To16() != nil && addr.IP.To4() == nil {
		family = syscall.AF_INET6
	}
	if addr.IP == nil || addr.IP.IsUnspecified() {
		family = syscall.AF_INET
	}
	println(network, address, family)
	s, err := socket(family, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return nil, err
	}
	sa, err := sockaddr(family, address)
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
