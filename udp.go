package muxedsocket

import (
	"net"
	"net/netip"
)

func init() {
	GlobalCreators().ChannelDialers().Register("udp", DialUDP)
	GlobalCreators().ChannelDialers().Register("udp", ListenUDP)
}

func DialUDP(addr string) (net.PacketConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	return net.DialUDP("udp", nil, udpAddr)
}
func ListenUDP(addr string) (net.PacketConn, error) {
	udpAddr, err := netip.ParseAddrPort(addr)
	if err != nil {
		return nil, err
	}
	return net.ListenUDP("udp", net.UDPAddrFromAddrPort(udpAddr))
}
