// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bulk

import "unsafe"

func (iov *sysIovec) set(b []byte) {
	iov.Base = (*byte)(unsafe.Pointer(&b[0]))
	iov.Len = uint64(len(b))
}
