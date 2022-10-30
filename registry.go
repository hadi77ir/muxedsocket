package muxedsocket

import (
	"crypto/tls"
	"net"
)

// TODO: type naming needs some improvement.

////////////////////////////////////////////////////////
// Client-side types

// PacketChannelDialer is used for example by UDP, ICMPChannel, etc.
type PacketChannelDialer func(addr string) (net.PacketConn, error)

// ClientStreamAdapter is used by packet-to-stream connection adapters. for example by KCP, eNet, etc.
type ClientStreamAdapter func(conn net.PacketConn, remoteAddr string) (net.Conn, error)

// ClientStreamDialer is used by packet-to-stream connection adapters. for example by KCP, eNet, etc.
type ClientStreamDialer func(addr string) (net.Conn, error)

// ClientMuxer is used for example by smux, yamux, SPDY, etc.
type ClientMuxer func(net.Conn, *ClientParams) (MuxedSocket, error)

// MuxDialer is the all-in-one solution: It takes a packet connection, returns a multiplexed secure streaming socket.
// Its sole purpose is for QUIC. (the only all-in-one solution we have at the moment!)
type MuxDialer func(conn net.PacketConn, config *tls.Config, params *ClientParams) (MuxedSocket, error)

// AddrMuxDialer is another all-in-one solution: It takes an address, returns the muxed socket connection.
type AddrMuxDialer func(addr string, config *tls.Config) (MuxedSocket, error)

////////////////////////////////////////////////////////
// Server-side types

// PacketChannelListener is used for example by UDP, ICMPChannel, etc.
type PacketChannelListener func(addr string) (net.PacketConn, error)

// ServerStreamAdapter is used by packet-to-stream connection adapters. for example by KCP, eNet, etc.
type ServerStreamAdapter func(conn net.PacketConn) (net.Listener, error)

// ServerMuxer is used for example by smux, yamux, SPDY, etc.
type ServerMuxer func(net.Listener, *ServerParams) (MuxedListener, error)

// MuxListener is the all-in-one solution: It takes a packet connection, returns a multiplexed secure streaming socket.
// Its sole purpose is for QUIC. (the only all-in-one solution we have at the moment!)
type MuxListener func(conn net.PacketConn, config *tls.Config, params *ServerParams) (MuxedListener, error)

// AddrMuxListener is another all-in-one solution, with the difference being: it takes an address instead of PacketConn.
type AddrMuxListener func(addr string, config *tls.Config) (MuxedListener, error)

////////////////////////////////////////////////////////
// We should have maps and registrar methods, per type.

type Registry[T any] struct {
	_initialized bool
	creators     map[string]T
}

func (r *Registry[T]) init() {
	r.creators = make(map[string]T)
	r._initialized = true
}

func (r *Registry[T]) Register(name string, value T) {
	if !r._initialized {
		r.init()
	}
	r.creators[name] = value
}
func (r *Registry[T]) Get(name string) (T, bool) {
	if !r._initialized {
		r.init()
	}
	if value, ok := r.creators[name]; ok {
		return value, true
	}
	var none T
	return none, false
}

// Client-side maps

type Creators struct {
	channelDialers       *Registry[PacketChannelDialer]
	clientStreamAdapters *Registry[ClientStreamAdapter]
	clientMuxers         *Registry[ClientMuxer]
	muxDialers           *Registry[MuxDialer]
	addrMuxDialers       *Registry[AddrMuxDialer]
	channelListeners     *Registry[PacketChannelListener]
	serverStreamAdapters *Registry[ServerStreamAdapter]
	serverMuxers         *Registry[ServerMuxer]
	muxListeners         *Registry[MuxListener]
	addrMuxListeners     *Registry[AddrMuxListener]
}

func (c *Creators) ChannelDialers() *Registry[PacketChannelDialer] {
	return c.channelDialers
}

func (c *Creators) ClientStreamAdapters() *Registry[ClientStreamAdapter] {
	return c.clientStreamAdapters
}

func (c *Creators) ClientMuxers() *Registry[ClientMuxer] {
	return c.clientMuxers
}

func (c *Creators) MuxDialers() *Registry[MuxDialer] {
	return c.muxDialers
}

func (c *Creators) AddrMuxDialers() *Registry[AddrMuxDialer] {
	return c.addrMuxDialers
}

func (c *Creators) ChannelListeners() *Registry[PacketChannelListener] {
	return c.channelListeners
}

func (c *Creators) ServerStreamAdapters() *Registry[ServerStreamAdapter] {
	return c.serverStreamAdapters
}

func (c *Creators) ServerMuxers() *Registry[ServerMuxer] {
	return c.serverMuxers
}

func (c *Creators) MuxListeners() *Registry[MuxListener] {
	return c.muxListeners
}

func (c *Creators) AddrMuxListeners() *Registry[AddrMuxListener] {
	return c.addrMuxListeners
}

var creators = NewCreators()

func GlobalCreators() *Creators {
	return creators
}
func NewCreators() *Creators {
	return &Creators{
		channelDialers:       &Registry[PacketChannelDialer]{},
		clientStreamAdapters: &Registry[ClientStreamAdapter]{},
		clientMuxers:         &Registry[ClientMuxer]{},
		muxDialers:           &Registry[MuxDialer]{},
		addrMuxDialers:       &Registry[AddrMuxDialer]{},
		channelListeners:     &Registry[PacketChannelListener]{},
		serverStreamAdapters: &Registry[ServerStreamAdapter]{},
		serverMuxers:         &Registry[ServerMuxer]{},
		muxListeners:         &Registry[MuxListener]{},
		addrMuxListeners:     &Registry[AddrMuxListener]{},
	}
}
