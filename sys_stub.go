// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux,!netbsd

package bulk

import "net"

const sysMSG_WAITFORONE = 0

type sysMmsghdr struct{}

func (msgs messages) scatter() []sysMmsghdr                      { return nil }
func (msgs *messages) gather(mmsgs []sysMmsghdr, laddr net.Addr) {}

func recvmmsg(s uintptr, b *Batch, flags uint32) (int, error) { return 0, errOpNoSupport }
func sendmmsg(s uintptr, b *Batch, flags uint32) (int, error) { return 0, errOpNoSupport }
