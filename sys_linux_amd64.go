// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bulk

import (
	"net"
	"syscall"
	"unsafe"
)

func (iov *sysIovec) set(b []byte) {
	iov.Base = (*byte)(unsafe.Pointer(&b[0]))
	iov.Len = uint64(len(b))
}

const (
	sysRECVMMSG = 299
	sysSENDMMSG = 307
)

func getsockname(s uintptr) (net.IP, int, string, error) {
	b := make([]byte, 128) // sizeofSockaddrStorage
	l := uint32(128)
	_, _, errno := syscall.RawSyscall(syscall.SYS_GETSOCKNAME, s, uintptr(unsafe.Pointer(&b[0])), uintptr(unsafe.Pointer(&l)))
	if errno != 0 {
		return nil, 0, "", error(errno)
	}
	ip, port, zone := ipPortZone((*byte)(unsafe.Pointer(&b[0])), l)
	return ip, port, zone, nil
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

func soclose(s uintptr) error { return syscall.Close(int(s)) }
