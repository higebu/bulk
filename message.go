// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bulk

import "net"

// A Message represents a datagaram message.
type Message struct {
	// N is a number of bytes received in read operation, a number
	// of bytes transferred in write operation.
	N int

	// Addr is a source address in read operation, a destination
	// address in write operation.
	Addr net.Addr

	Data []byte
}

type messages []Message

// A Batch represents a batch of messages.
type Batch struct {
	Messages []Message

	msgs []sysMmsghdr
}

// Scatter scatters b on internal data for the following bacth read,
// write operations.
func (b *Batch) Scatter() error {
	b.msgs = messages(b.Messages).scatter()
	return nil
}

// Reset resets b.
// At present, it just sets each N of Messages in b to zero.
func (b *Batch) Reset() {
	for i := range b.Messages {
		b.Messages[i].N = 0
	}
}
