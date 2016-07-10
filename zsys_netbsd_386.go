// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs defs_netbsd.go

package bulk

const (
	sysMSG_DONTWAIT   = 0x80
	sysMSG_WAITFORONE = 0x2000

	sizeofSockaddrInet  = 0x10
	sizeofSockaddrInet6 = 0x1c
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
	Hdr sysMsghdr
	Len uint32
}

type sysMsghdr struct {
	Name       *byte
	Namelen    uint32
	Iov        *sysIovec
	Iovlen     int32
	Control    *byte
	Controllen uint32
	Flags      int32
}

type sysIovec struct {
	Base *byte
	Len  uint32
}

type sysTimespec struct {
	Sec  int64
	Nsec int32
}
