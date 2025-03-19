package types

import (
	"net"
	"net/netip"
	"strconv"
	"strings"
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

func GetAddrByTransportType(transportKey string, addr string) (NetAddrPort, error) {
	sep := strings.LastIndex(addr, ":")
	var ipStr string
	var portStr string
	if sep != -1 {
		ipStr = addr[:sep]
		portStr = addr[sep+1:]
	} else {
		ipStr = addr
		portStr = "0"
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, net.InvalidAddrError("invalid ip address")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	if transportKey == "udp" {
		return &net.UDPAddr{IP: ip, Port: port}, nil
	}
	if transportKey == "tcp" {
		return &net.TCPAddr{IP: ip, Port: port}, nil
	}

	return &netAddr{
		host:    addr,
		port:    uint16(port),
		network: transportKey,
	}, nil
}

type EmptyAddr string

func (e EmptyAddr) Network() string {
	return string(e)
}

func (e EmptyAddr) String() string {
	return string(e) + ":"
}

func (e EmptyAddr) AddrPort() netip.AddrPort {
	return netip.MustParseAddrPort("0.0.0.0:0")
}

var _ NetAddrPort = EmptyAddr("empty")
