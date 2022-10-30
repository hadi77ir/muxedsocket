package muxedsocket

import (
	"net"
	"time"
)

type MuxedListener interface {
	ListeningSocket
}

type ListeningSocket interface {
	// Accept waits for and returns the next connection to the listener.
	Accept() (socket Socket, err error)

	// Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() error

	// Addr returns the listener's network address.
	Addr() net.Addr
}

type Socket interface {
	// Close closes the connection.
	// Any blocked operations on the socket will be unblocked and return errors.
	Close() error

	// LocalAddr returns the local network address, if known.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address, if known.
	RemoteAddr() net.Addr
}

type MuxedSocket interface {
	Socket
	AcceptStream() (stream MuxStream, err error)
	OpenStream() (stream MuxStream, err error)
}

type DatagramChannel interface {
	SendDatagram(b []byte) error
	ReceiveDatagram() ([]byte, error)
}

type DatagramCapableMuxedSocket interface {
	MuxedSocket
	DatagramChannel
}

type MuxStream interface {
	Socket
	net.Conn
	StreamID() int
}

type CommonParams struct {
	KeepalivePeriod time.Duration
	MaxIdleTimeout  time.Duration
}

type ClientParams struct {
	CommonParams
	// used with "WriteTo" method. a PacketConn only has "LocalAddr" of the local socket.
	RemoteAddr net.Addr
}

type ServerParams struct {
	CommonParams
}
