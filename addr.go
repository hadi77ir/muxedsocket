package muxedsocket

import (
	"net"
	"net/netip"
	"strconv"
)

type MuxedAddr struct {
	Original  net.Addr
	MuxStream int
}

// String returns a serialized string representation of struct.
func (addr MuxedAddr) String() string {
	return addr.Original.String() + "#" + strconv.Itoa(addr.MuxStream)
}

// Network returns network name, such as "icmp", "tcp", "udp".
func (addr MuxedAddr) Network() string {
	return addr.Original.Network()
}

func WrapAddr(addr net.Addr, id int) *MuxedAddr {
	return &MuxedAddr{Original: addr, MuxStream: id}
}

type NetAddrPort interface {
	net.Addr
	AddrPort() netip.AddrPort
}

type netAddr struct {
	// host contains "host" or "host:port"
	host    string
	port    uint16
	network string
	zone    string
}

func (n netAddr) Network() string {
	return n.network
}

func (n netAddr) Port() uint16 {
	return n.port
}

func (n netAddr) AddrPort() netip.AddrPort {
	s, err := netip.ParseAddrPort(n.host)
	if err != nil {
		var none netip.AddrPort
		return none
	}
	return s
}

func (n netAddr) String() string {
	return n.host
}

func (n netAddr) Zone() string {
	return n.zone
}
