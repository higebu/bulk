// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bulk

import (
	"net"
	"syscall"
	"unsafe"
)

func (iov *iovec) set(b []byte) {
	iov.Base = (*byte)(unsafe.Pointer(&b[0]))
	iov.Len = uint32(len(b))
}

const (
	sysGETSOCKNAME = 6
	sysGETPEERNAME = 7
	sysRECVMMSG    = 19
	sysSENDMMSG    = 20
)

func socketcall(call int, a0, a1, a2, a3, a4, a5 uintptr) (int, syscall.Errno)

func getsockname(s uintptr) (net.IP, int, string, error) {
	b := make([]byte, 128) // sizeofSockaddrStorage
	l := uint32(128)
	_, errno := socketcall(sysGETSOCKNAME, s, uintptr(unsafe.Pointer(&b[0])), uintptr(unsafe.Pointer(&l)), 0, 0, 0)
	if errno != 0 {
		return nil, 0, "", error(errno)
	}
	ip, port, zone := ipPortZone((*byte)(unsafe.Pointer(&b[0])), l)
	return ip, port, zone, nil
}

func getpeername(s uintptr) (net.IP, int, string, error) {
	b := make([]byte, 128) // sizeofSockaddrStorage
	r := uint32(128)
	_, errno := socketcall(sysGETPEERNAME, s, uintptr(unsafe.Pointer(&b[0])), uintptr(unsafe.Pointer(&r)), 0, 0, 0)
	if errno != 0 {
		return nil, 0, "", error(errno)
	}
	ip, port, zone := ipPortZone((*byte)(unsafe.Pointer(&b[0])), r)
	return ip, port, zone, nil
}

func recvmmsg(s uintptr, mmsgs []mmsghdr, flags uint32) (int, error) {
	l := uint32(len(mmsgs))
	n, errno := socketcall(sysRECVMMSG, s, uintptr(unsafe.Pointer(&mmsgs[0])), uintptr(l), uintptr(flags), 0, 0)
	if errno != 0 {
		return 0, error(errno)
	}
	return int(n), nil
}

func sendmmsg(s uintptr, mmsgs []mmsghdr, flags uint32) (int, error) {
	l := uint32(len(mmsgs))
	n, errno := socketcall(sysSENDMMSG, s, uintptr(unsafe.Pointer(&mmsgs[0])), uintptr(l), uintptr(flags), 0, 0)
	if errno != 0 {
		return 0, error(errno)
	}
	return int(n), nil
}
