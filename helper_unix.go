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
		if ip = ip.To4(); ip == nil {
			return nil, net.InvalidAddrError("non-ipv4 address")
		}
		sa := &syscall.SockaddrInet4{Port: port}
		copy(sa.Addr[:], ip)
		return sa, nil
	case syscall.AF_INET6:
		if len(ip) == 0 {
			ip = net.IPv6unspecified
		}
		if ip.Equal(net.IPv4zero) {
			ip = net.IPv6unspecified
		}
		if ip = ip.To16(); ip == nil || ip.To4() != nil {
			return nil, net.InvalidAddrError("non-ipv6 address")
		}
		sa := &syscall.SockaddrInet6{Port: port, ZoneId: uint32(zoneCache.nameToIndex(zone))}
		copy(sa.Addr[:], ip)
		return sa, nil
	default:
		return nil, net.InvalidAddrError("unexpected family")
	}
}
