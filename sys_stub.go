// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux,!netbsd linux,!amd64,!386 netbsd,!amd64,!386

package bulk

import (
	"errors"
	"net"
)

const sysMSG_WAITFORONE = 0

type mmsghdr struct{}

func (msgs messages) scatter() []mmsghdr                      { return nil }
func (msgs *messages) gather(mmsgs []mmsghdr, laddr net.Addr) {}

func recvmmsg(s uintptr, mmsgs []mmsghdr, flags uint32) (int, error) {
	return 0, errors.New("operation not supported")
}

func sendmmsg(s uintptr, mmsgs []mmsghdr, flags uint32) (int, error) {
	return 0, errors.New("operation not supported")
}

func socket(family, sotype, proto int) (uintptr, error) {
	return 0, errors.New("operation not supported")
}

func getsockname(s uintptr) (net.IP, int, string, error) {
	return nil, 0, "", errors.New("operation not supported")
}

func getpeername(s uintptr) (net.IP, int, string, error) {
	return nil, 0, "", errors.New("operation not supported")
}

func soclose(s uintptr) error {
	return errors.New("operation not supported")
}
