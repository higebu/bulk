// Copyright 2016 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// +godefs map struct_in_addr [4]byte /* in_addr */
// +godefs map struct_in6_addr [16]byte /* in6_addr */

package bulk

/*
#include <sys/types.h>
#define __USE_GNU
#include <sys/socket.h>

#include <linux/in.h>
#include <linux/in6.h>
*/
import "C"

const (
	sysMSG_DONTWAIT   = C.MSG_DONTWAIT
	sysMSG_WAITFORONE = C.MSG_WAITFORONE

	sizeofSockaddrInet  = C.sizeof_struct_sockaddr_in
	sizeofSockaddrInet6 = C.sizeof_struct_sockaddr_in6
)

type sockaddrInet C.struct_sockaddr_in

type sockaddrInet6 C.struct_sockaddr_in6

type mmsghdr C.struct_mmsghdr

type msghdr C.struct_msghdr

type iovec C.struct_iovec

type timespec C.struct_timespec
