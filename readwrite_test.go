// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux netbsd

package bulk_test

import (
	"net"
	"sync"
	"testing"

	"github.com/higebu/bulk"
)

func TestReadWrite(t *testing.T) {
	c, err := bulk.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
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

func TestReadWriteWithDial(t *testing.T) {
	s, err := bulk.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	laddr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	c, err := bulk.DialPacket("udp4", laddr, s.LocalAddr().(*net.UDPAddr))
	if err != nil {
		t.Fatal(err)
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
		wb.Messages[i].Addr = s.LocalAddr()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	rb.Scatter()
	go func() {
		defer wg.Done()
		var nr int
		for nr < N {
			n, err := s.ReadBatch(&rb)
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
