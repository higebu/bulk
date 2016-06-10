// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux,!netbsd linux,!amd64,!386 netbsd,!amd64

package bulk

import "net"

const sysMSG_WAITFORONE = 0

type sysMmsghdr struct{}

func (msgs messages) scatter() []sysMmsghdr                      { return nil }
func (msgs *messages) gather(mmsgs []sysMmsghdr, laddr net.Addr) {}

func recvmmsg(s uintptr, mmsgs []sysMmsghdr, flags uint32) (int, error) { return 0, errOpNoSupport }
func sendmmsg(s uintptr, mmsgs []sysMmsghdr, flags uint32) (int, error) { return 0, errOpNoSupport }

func socket(family, sotype, proto int) (uintptr, error)  { return 0, errOpNoSupport }
func getsockname(s uintptr) (net.IP, int, string, error) { return nil, 0, "", errOpNoSupport }
func soclose(s uintptr) error                            { return errOpNoSupport }
