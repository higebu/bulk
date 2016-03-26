// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,amd64 linux,386 netbsd,amd64

package bulk

import (
	"net"
	"unsafe"
)

func (sa *sysSockaddrInet) ipPortZone() (net.IP, int, string) {
	ip := net.IPv4(sa.Addr[0], sa.Addr[1], sa.Addr[2], sa.Addr[3])
	p := (*[2]byte)(unsafe.Pointer(&sa.Port))
	port := int(p[0])<<8 + int(p[1])
	return ip, port, ""
}

func (sa *sysSockaddrInet6) ipPortZone() (net.IP, int, string) {
	ip := make([]byte, net.IPv6len)
	copy(ip, sa.Addr[:])
	p := (*[2]byte)(unsafe.Pointer(&sa.Port))
	port := int(p[0])<<8 + int(p[1])
	return ip, port, ifIB.indexToName(int(sa.Scope_id))
}

func ipPortZone(b *byte, l uint32) (net.IP, int, string) {
	if l == sysSizeofSockaddrInet {
		sa := (*sysSockaddrInet)(unsafe.Pointer(b))
		return sa.ipPortZone()
	}
	if l == sysSizeofSockaddrInet6 {
		sa := (*sysSockaddrInet6)(unsafe.Pointer(b))
		return sa.ipPortZone()
	}
	return nil, 0, ""
}

func (msgs messages) scatter() []sysMmsghdr {
	mmsgs := make([]sysMmsghdr, 0, len(msgs))
	for i := range msgs {
		if len(msgs[i].Data) == 0 {
			continue
		}
		var m sysMmsghdr
		switch addr := msgs[i].Addr.(type) {
		case *net.UDPAddr:
			m.Hdr.Name, m.Hdr.Namelen = msgSockaddr(addr.IP, addr.Port, addr.Zone)
		case *net.IPAddr:
			m.Hdr.Name, m.Hdr.Namelen = msgSockaddr(addr.IP, 0, addr.Zone)
		}
		iov := &sysIovec{}
		iov.set(msgs[i].Data)
		m.Hdr.Iov = iov
		m.Hdr.Iovlen = 1
		mmsgs = append(mmsgs, m)
	}
	return mmsgs
}

func (msgs *messages) gather(mmsgs []sysMmsghdr, laddr net.Addr) {
	for i := range mmsgs {
		var addr net.Addr
		switch laddr.(type) {
		case *net.UDPAddr:
			udp := &net.UDPAddr{}
			udp.IP, udp.Port, udp.Zone = ipPortZone(mmsgs[i].Hdr.Name, mmsgs[i].Hdr.Namelen)
			addr = udp
		case *net.IPAddr:
			ip := &net.IPAddr{}
			ip.IP, _, ip.Zone = ipPortZone(mmsgs[i].Hdr.Name, mmsgs[i].Hdr.Namelen)
			addr = ip
		default:
		}
		if addr != nil {
			(*msgs)[i].Addr = addr
		}
		(*msgs)[i].N = int(mmsgs[i].Len)
	}
}
