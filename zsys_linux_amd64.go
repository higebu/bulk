// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs defs_linux.go

package bulk

const (
	sysMSG_DONTWAIT   = 0x40
	sysMSG_WAITFORONE = 0x10000

	sizeofSockaddrInet  = 0x10
	sizeofSockaddrInet6 = 0x1c
)

type sockaddrInet struct {
	Family uint16
	Port   uint16
	Addr   [4]byte /* in_addr */
	X__pad [8]uint8
}

type sockaddrInet6 struct {
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
}

type mmsghdr struct {
	Hdr       msghdr
	Len       uint32
	Pad_cgo_0 [4]byte
}

type msghdr struct {
	Name       *byte
	Namelen    uint32
	Pad_cgo_0  [4]byte
	Iov        *iovec
	Iovlen     uint64
	Control    *byte
	Controllen uint64
	Flags      int32
	Pad_cgo_1  [4]byte
}

type iovec struct {
	Base *byte
	Len  uint64
}

type timespec struct {
	Sec  int64
	Nsec int64
}
