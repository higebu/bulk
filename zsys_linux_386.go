// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs defs_linux.go

package bulk

const (
	sysMSG_DONTWAIT   = 0x40
	sysMSG_WAITFORONE = 0x10000

	sysSizeofSockaddrInet  = 0x10
	sysSizeofSockaddrInet6 = 0x1c
)

type sysSockaddrInet struct {
	Family uint16
	Port   uint16
	Addr   [4]byte /* in_addr */
	X__pad [8]uint8
}

type sysSockaddrInet6 struct {
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
}

type sysMmsghdr struct {
	Hdr sysMsghdr
	Len uint32
}

type sysMsghdr struct {
	Name       *byte
	Namelen    uint32
	Iov        *sysIovec
	Iovlen     uint32
	Control    *byte
	Controllen uint32
	Flags      int32
}

type sysIovec struct {
	Base *byte
	Len  uint32
}

type sysTimespec struct {
	Sec  int32
	Nsec int32
}