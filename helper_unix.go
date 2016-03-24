// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux netbsd

package bulk

import (
	"net"
	"strconv"
	"syscall"
)

func sockaddr(family int, address string) (syscall.Sockaddr, error) {
	switch family {
	case syscall.AF_INET:
		a, err := net.ResolveUDPAddr("udp4", address)
		if err != nil {
			return nil, err
		}
		if len(a.IP) == 0 {
			a.IP = net.IPv4zero
		}
		if a.IP = a.IP.To4(); a.IP == nil {
			return nil, net.InvalidAddrError("non-ipv4 address")
		}
		sa := &syscall.SockaddrInet4{Port: a.Port}
		copy(sa.Addr[:], a.IP)
		return sa, nil
	case syscall.AF_INET6:
		a, err := net.ResolveUDPAddr("udp6", address)
		if err != nil {
			return nil, err
		}
		if len(a.IP) == 0 {
			a.IP = net.IPv6unspecified
		}
		if a.IP.Equal(net.IPv4zero) {
			a.IP = net.IPv6unspecified
		}
		if a.IP = a.IP.To16(); a.IP == nil || a.IP.To4() != nil {
			return nil, net.InvalidAddrError("non-ipv6 address")
		}
		sa := &syscall.SockaddrInet6{Port: a.Port, ZoneId: zoneToUint32(a.Zone)}
		copy(sa.Addr[:], a.IP)
		return sa, nil
	default:
		return nil, net.InvalidAddrError("unexpected family")
	}
}

func zoneToUint32(zone string) uint32 {
	if zone == "" {
		return 0
	}
	if ifi, err := net.InterfaceByName(zone); err == nil {
		return uint32(ifi.Index)
	}
	n, err := strconv.Atoi(zone)
	if err != nil {
		return 0
	}
	return uint32(n)
}
