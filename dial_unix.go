// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,amd64 linux,386 netbsd,amd64 netbsd,386

package bulk

import (
	"net"
	"os"
	"syscall"
)

func dialPacket(network string, laddr, raddr *net.UDPAddr) (*PacketConn, error) {
	var family int
	if laddr != nil {
		family = addrFamily(network, laddr.IP)
	} else {
		family = addrFamily(network, nil)
	}
	s, err := socket(family, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		return nil, err
	}

	if laddr != nil {
		lsa, err := sockaddr(family, laddr.IP, laddr.Port, laddr.Zone)
		if err != nil {
			syscall.Close(int(s))
			return nil, err
		}
		if err := syscall.Bind(int(s), lsa); err != nil {
			syscall.Close(int(s))
			return nil, os.NewSyscallError("bind", err)
		}
	}

	rsa, err := sockaddr(family, raddr.IP, raddr.Port, raddr.Zone)
	if err != nil {
		syscall.Close(int(s))
		return nil, err
	}
	if err := syscall.Connect(int(s), rsa); err != nil {
		syscall.Close(int(s))
		return nil, os.NewSyscallError("connect", err)
	}

	lip, lport, lzone, err := getsockname(s)
	if err != nil {
		syscall.Close(int(s))
		return nil, os.NewSyscallError("getsockname", err)
	}
	rip, rport, rzone, err := getpeername(s)
	if err != nil {
		syscall.Close(int(s))
		return nil, os.NewSyscallError("getsockname", err)
	}
	return &PacketConn{
		s:     s,
		laddr: &net.UDPAddr{IP: lip, Port: lport, Zone: lzone},
		raddr: &net.UDPAddr{IP: rip, Port: rport, Zone: rzone},
	}, nil
}
