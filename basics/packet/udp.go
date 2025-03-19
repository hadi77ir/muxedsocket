package packet

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
	"net"
	"net/netip"
)

func NewUDPImplementation() types.PacketConnImplementation {
	return WrapImplementation("udp", dialUDP, listenUDP, resolveUDP)
}

func dialUDP(network string, addr string) (net.PacketConn, error) {
	resolved, err := resolveUDP(network, addr)
	if err != nil {
		return nil, err
	}
	uAddr, ok := resolved.(*net.UDPAddr)
	if !ok {
		return nil, muxedsocket.ErrOpNotSupported
	}
	return net.DialUDP("udp", nil, uAddr)
}

func listenUDP(network string, addr string) (net.PacketConn, error) {
	resolved, err := resolveUDP(network, addr)
	if err != nil {
		return nil, err
	}
	uAddr, ok := resolved.(*net.UDPAddr)
	if !ok {
		return nil, muxedsocket.ErrOpNotSupported
	}
	return net.ListenUDP(network, uAddr)
}
func resolveUDP(network, addr string) (net.Addr, error) {
	addrPort, err := netip.ParseAddrPort(addr)
	if err != nil {
		resolved, err := net.ResolveUDPAddr(network, addr)
		if err != nil {
			return nil, err
		}
		return resolved, nil
	}
	return net.UDPAddrFromAddrPort(addrPort), nil
}
