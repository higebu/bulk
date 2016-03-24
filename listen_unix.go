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
	var family, proto int
	switch network {
	case "udp4":
		family, proto = syscall.AF_INET, syscall.IPPROTO_UDP
	case "udp6":
		family, proto = syscall.AF_INET6, syscall.IPPROTO_UDP
	default:
		return nil, errOpNoSupport
	}
	s, err := syscall.Socket(family, syscall.SOCK_DGRAM, proto)
	if err != nil {
		return nil, os.NewSyscallError("socket", err)
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
	ip, port, err := getsockname(uintptr(s))
	if err != nil {
		syscall.Close(s)
		return nil, os.NewSyscallError("getsockname", err)
	}
	c := PacketConn{s: s}
	c.laddr = &net.UDPAddr{IP: ip, Port: port}
	return &c, nil
}
