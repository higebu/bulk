// Copyright 2017 Yuya Kusakabe. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux,!netbsd linux,!amd64,!386 netbsd,!amd64,!386

package bulk

import (
	"errors"
	"net"
)

func dialPacket(network, laddr, raddr *net.UDPAddr) (*PacketConn, error) {
	return nil, errors.New("operation not supported")
}
