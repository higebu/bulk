// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux netbsd

package bulk_test

import (
	"testing"

	"github.com/mikioh/bulk"
)

func BenchmarkReadWrite(b *testing.B) {
	c, err := bulk.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	defer c.Close()

	N := 16            // see UIO_MAXIOV or similar caps
	const dataLen = 18 // UDP over IPv4 over Ethernet

	var rb, wb bulk.Batch
	rb.Messages = make([]bulk.Message, N)
	wb.Messages = make([]bulk.Message, N)
	for i := 0; i < N; i++ {
		rb.Messages[i].Data = make([]byte, dataLen)
		wb.Messages[i].Data = make([]byte, dataLen)
		wb.Messages[i].Addr = c.LocalAddr()
	}

	go func() {
		wb.Scatter()
		for {
			if _, err := c.WriteBatch(&wb); err != nil {
				break
			}
		}
	}()

	rb.Scatter()
	for i := 0; i < b.N; i++ {
		if _, err := c.ReadBatch(&rb); err != nil {
			b.Fatal(err)
		}
		rb.Reset()
	}
}
