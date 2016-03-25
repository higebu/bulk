// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !linux,!netbsd !amd64

package bulk

func listenPacket(network, address string) (*PacketConn, error) { return nil, errOpNoSupport }
