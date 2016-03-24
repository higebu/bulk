// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bulk_test

import (
	"runtime"
	"sync"
	"testing"

	"github.com/mikioh/bulk"
)

func TestReadWrite(t *testing.T) {
	c, err := bulk.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	N := 16 // see UIO_MAXIOV or similar caps
	if runtime.GOOS == "linux" {
		N = 256
	}
	const dataLen = 18 // UDP over IPv4 over Ethernet

	var rb, wb bulk.Batch
	rb.Messages = make([]bulk.Message, N)
	wb.Messages = make([]bulk.Message, N)
	for i := 0; i < N; i++ {
		rb.Messages[i].Data = make([]byte, dataLen)
		wb.Messages[i].Data = make([]byte, dataLen)
		wb.Messages[i].Addr = c.LocalAddr()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	rb.Scatter()
	go func() {
		defer wg.Done()
		var nr int
		for nr < N {
			n, err := c.ReadBatch(&rb)
			if err != nil {
				t.Error(err)
				return
			}
			nr += n
			for _, m := range rb.Messages[:n] {
				if m.Addr == nil {
					t.Errorf("got %#v; want non-nil", m.Addr)
				}
			}
		}
	}()
	wg.Add(1)
	wb.Scatter()
	go func() {
		defer wg.Done()
		nw, err := c.WriteBatch(&wb)
		if err != nil {
			t.Error(err)
			return
		}
		if nw != N {
			t.Errorf("got %d; want %d", nw, N)
			return
		}
	}()
	wg.Wait()
}
