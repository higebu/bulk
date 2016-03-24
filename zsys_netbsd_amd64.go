// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs defs_netbsd.go

package bulk

const (
	sysMSG_DONTWAIT   = 0x80
	sysMSG_WAITFORONE = 0x2000

	sysSizeofSockaddrInet  = 0x10
	sysSizeofSockaddrInet6 = 0x1c
)

type sysSockaddrInet struct {
	Len    uint8
	Family uint8
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]int8
}

type sysSockaddrInet6 struct {
	Len      uint8
	Family   uint8
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
}

type sysMmsghdr struct {
	Hdr       sysMsghdr
	Len       uint32
	Pad_cgo_0 [4]byte
}

type sysMsghdr struct {
	Name       *byte
	Namelen    uint32
	Pad_cgo_0  [4]byte
	Iov        *sysIovec
	Iovlen     int32
	Pad_cgo_1  [4]byte
	Control    *byte
	Controllen uint32
	Flags      int32
}

type sysIovec struct {
	Base *byte
	Len  uint64
}

type sysTimespec struct {
	Sec  int64
	Nsec int64
}
