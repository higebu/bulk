// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,amd64 netbsd,amd64

package bulk

import (
	"net"
	"syscall"
	"unsafe"
)

func getsockname(s uintptr) (net.IP, int, string, error) {
	b := make([]byte, 128) // sysSizeofSockaddrStorage
	l := uint32(128)
	_, _, errno := syscall.RawSyscall(syscall.SYS_GETSOCKNAME, s, uintptr(unsafe.Pointer(&b[0])), uintptr(unsafe.Pointer(&l)))
	if errno != 0 {
		return nil, 0, "", error(errno)
	}
	ip, port, zone := ipPortZone((*byte)(unsafe.Pointer(&b[0])), l)
	return ip, port, zone, nil
}

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
	return ip, port, zoneToString(sa.Scope_id)
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
		var iov sysIovec
		iov.Base = (*byte)(unsafe.Pointer(&msgs[i].Data[0]))
		iov.Len = uint64(len(msgs[i].Data))
		m.Hdr.Iov = &iov
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

func recvmmsg(s uintptr, mmsgs []sysMmsghdr, flags uint32) (int, error) {
	l := uint32(len(mmsgs))
	n, _, errno := syscall.Syscall6(sysRECVMMSG, s, uintptr(unsafe.Pointer(&mmsgs[0])), uintptr(l), uintptr(flags), 0, 0)
	if errno != 0 {
		return 0, error(errno)
	}
	return int(n), nil
}

func sendmmsg(s uintptr, mmsgs []sysMmsghdr, flags uint32) (int, error) {
	l := uint32(len(mmsgs))
	n, _, errno := syscall.Syscall6(sysSENDMMSG, s, uintptr(unsafe.Pointer(&mmsgs[0])), uintptr(l), uintptr(flags), 0, 0)
	if errno != 0 {
		return 0, error(errno)
	}
	return int(n), nil
}
