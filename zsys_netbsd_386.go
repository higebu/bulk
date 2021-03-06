// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs defs_netbsd.go

package bulk

const (
	sysMSG_DONTWAIT   = 0x80
	sysMSG_WAITFORONE = 0x2000

	sizeofSockaddrInet  = 0x10
	sizeofSockaddrInet6 = 0x1c
)

type sockaddrInet struct {
	Len    uint8
	Family uint8
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]int8
}

type sockaddrInet6 struct {
	Len      uint8
	Family   uint8
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
}

type mmsghdr struct {
	Hdr msghdr
	Len uint32
}

type msghdr struct {
	Name       *byte
	Namelen    uint32
	Iov        *iovec
	Iovlen     int32
	Control    *byte
	Controllen uint32
	Flags      int32
}

type iovec struct {
	Base *byte
	Len  uint32
}

type timespec struct {
	Sec  int64
	Nsec int32
}
