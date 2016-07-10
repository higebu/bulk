// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux netbsd

package bulk

import (
	"net"
	"syscall"
)

func addrFamily(network string, ip net.IP) int {
	if ip.To4() != nil {
		return syscall.AF_INET
	}
	if ip.To16() != nil && ip.To4() == nil {
		return syscall.AF_INET6
	}
	if ip == nil || ip.IsUnspecified() {
		if network == "udp4" {
			return syscall.AF_INET
		}
		if network == "udp6" {
			return syscall.AF_INET6
		}
	}
	return syscall.AF_INET
}

func sockaddr(family int, ip net.IP, port int, zone string) (syscall.Sockaddr, error) {
	switch family {
	case syscall.AF_INET:
		if len(ip) == 0 {
			ip = net.IPv4zero
		}
		ip4 := ip.To4()
		if ip4 == nil {
			return nil, &net.AddrError{Err: "non-IPv4 address", Addr: ip.String()}
		}
		sa := &syscall.SockaddrInet4{Port: port}
		copy(sa.Addr[:], ip4)
		return sa, nil
	case syscall.AF_INET6:
		if len(ip) == 0 || ip.Equal(net.IPv4zero) {
			ip = net.IPv6unspecified
		}
		ip6 := ip.To16()
		if ip6 == nil {
			return nil, &net.AddrError{Err: "non-IPv6 address", Addr: ip.String()}
		}
		sa := &syscall.SockaddrInet6{Port: port, ZoneId: uint32(zoneCache.index(zone))}
		copy(sa.Addr[:], ip6)
		return sa, nil
	default:
		return nil, &net.AddrError{Err: "invalid address family", Addr: ip.String()}
	}
}
