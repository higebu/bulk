// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bulk

import (
	"errors"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
)

var (
	_ net.PacketConn = &PacketConn{}

	errClosing     = errors.New("use of closed network connection")
	errOpNoSupport = errors.New("operation not supported")
)

// A PacketConn represents a packet network endpoint.
type PacketConn struct {
	mu    sync.RWMutex
	s     int
	laddr net.Addr
}

// ReadFrom reads a message from the endpoint.
func (c *PacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	return 0, nil, errOpNoSupport
}

// WriteTo writes the message b to dst.
func (c *PacketConn) WriteTo(b []byte, dst net.Addr) (int, error) {
	return 0, errOpNoSupport
}

// Close closes the endpoint.
func (c *PacketConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.s < 0 {
		return errClosing
	}
	err := syscall.Close(int(c.s))
	c.s = -1
	if err != nil {
		return os.NewSyscallError("close", err)
	}
	return nil
}

// LocalAddr returns the local network address.
func (c *PacketConn) LocalAddr() net.Addr {
	return c.laddr
}

// SetDeadline sets the read and write deadlines associated with the
// endpoint.
func (c *PacketConn) SetDeadline(t time.Time) error {
	return errOpNoSupport
}

// SetReadDeadline sets the read deadline associated with the
// endpoint.
func (c *PacketConn) SetReadDeadline(t time.Time) error {
	return errOpNoSupport
}

// SetWriteDeadline sets the write deadline associated with the
// endpoint.
func (c *PacketConn) SetWriteDeadline(t time.Time) error {
	return errOpNoSupport
}

// ReadBatch reads a batch of messages.
// The result of operation will be stored in b when no error occurs.
// The batch b must be scatterd before operation.
func (c *PacketConn) ReadBatch(b *Batch) (int, error) {
	var s int
	c.mu.RLock()
	if c.s < 0 {
		c.mu.RUnlock()
		return 0, errClosing
	}
	s = c.s
	c.mu.RUnlock()
	n, err := recvmmsg(uintptr(s), b.mmsgs, sysMSG_WAITFORONE)
	if err != nil {
		return 0, err
	}
	msgs := messages(b.Messages)
	(&msgs).gather(b.mmsgs[:n], c.laddr)
	return n, nil
}

// WriteBatch writes a batch of messages.
// The result of operation will be stored in b when no error occurs.
// The batch b must be scatterd before operation.
func (c *PacketConn) WriteBatch(b *Batch) (int, error) {
	var s int
	c.mu.RLock()
	if c.s < 0 {
		c.mu.RUnlock()
		return 0, errClosing
	}
	s = c.s
	c.mu.RUnlock()
	n, err := sendmmsg(uintptr(s), b.mmsgs, 0)
	if err != nil {
		return 0, err
	}
	msgs := messages(b.Messages)
	(&msgs).gather(b.mmsgs[:n], c.laddr)
	return n, nil
}

// ListenPacket listens for incoming datagrams addressed to address.
// At present, network must be "udp", "udp4" or "udp6".
// See net.Dial for the syntax of address.
func ListenPacket(network, address string) (*PacketConn, error) {
	return listenPacket(network, address)
}
