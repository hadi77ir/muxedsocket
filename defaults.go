package muxedsocket

import "github.com/hadi77ir/muxedsocket/types"

const (
	// DefaultPacketConn is set to "udp" so by default when there is no basic connection layer defined,
	// UDP may be used as basic packet-oriented transport.
	DefaultPacketConn = "udp"
	// DefaultStreamConn is set to "tcp" so by default when there is no basic connection layer defined,
	// TCP may be used as basic streaming transport.
	DefaultStreamConn = "tcp"
	// DefaultStreamSolution is set to "nomux". If the final result is a streaming connection and a multiplexer is desired,
	// this will be added on top of result.
	DefaultStreamSolution = "nomux"
	// DefaultPacketSolution is set to "quic". If the final result is a packet connection and a multiplexer is desired,
	// this will be added on top of result.
	DefaultPacketSolution = "quic"
	// DefaultStreamAdapter is set to "kcp".
	DefaultStreamAdapter = "kcp"
	// DefaultPacketAdapter is set to "spos" (Session-based Packets-over-streams).
	DefaultPacketAdapter = "spos"
)

type DefaultLayers struct {
	PacketConn     types.PacketConnImplementation
	StreamConn     types.StreamConnImplementation
	PacketSolution types.PacketSolutionImplementation
	StreamSolution types.StreamSolutionImplementation
	PacketAdapter  types.PacketAdapterImplementation
	StreamAdapter  types.StreamAdapterImplementation
	AddrSolution   types.AddrSolutionImplementation
}

func GetDefaults(c *Creators) *DefaultLayers {
	defaultStreamConn, _ := c.StreamConns().Get(DefaultStreamConn)
	defaultPacketConn, _ := c.PacketConns().Get(DefaultPacketConn)
	defaultPacketSolution, _ := c.PacketSolutions().Get(DefaultPacketSolution)
	defaultStreamSolution, _ := c.StreamSolutions().Get(DefaultStreamSolution)
	defaultPacketAdapter, _ := c.PacketAdapters().Get(DefaultPacketAdapter)
	defaultStreamAdapter, _ := c.StreamAdapters().Get(DefaultStreamAdapter)
	return &DefaultLayers{
		StreamConn:     defaultStreamConn,
		PacketConn:     defaultPacketConn,
		StreamSolution: defaultStreamSolution,
		PacketSolution: defaultPacketSolution,
		PacketAdapter:  defaultPacketAdapter,
		StreamAdapter:  defaultStreamAdapter,
	}
}
