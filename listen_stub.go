// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux,!netbsd linux,!amd64,!386 netbsd,!amd64,!386

package bulk

import "errors"

func listenPacket(network, address string) (*PacketConn, error) {
	return nil, errors.New("operation not supported")
}
