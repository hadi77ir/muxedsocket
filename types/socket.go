package types

import (
	"net"
)

type MuxedListener interface {
	ListeningSocket
	// AcceptMuxed waits for and returns the next connection to the listener. Same as Accept.
	AcceptMuxed() (socket MuxedSocket, err error)
}

type Closable interface {
	// CloseChan returns a read-only channel that is closed when the socket is closed.
	CloseChan() <-chan struct{}

	// Close closes the connection.
	// Any blocked operations on the socket will be unblocked and return errors.
	Close() error
}

type ListeningSocket interface {
	// Closable means this listener may be closed. Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Closable

	// Accept waits for and returns the next connection to the listener.
	Accept() (socket Socket, err error)

	// Addr returns the listener's network address.
	Addr() net.Addr
}

type Socket interface {
	Closable
	// LocalAddr returns the local network address, if known.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address, if known.
	RemoteAddr() net.Addr

	// CanRedial returns true if this connection can be redialed. This will return false if the
	// connection was accepted by a listener.
	CanRedial() bool

	// Redial will try to redial the remote and return the established connection.
	// It will return false if the current connection was accepted by a listener.
	Redial() (Socket, error)
}
type StreamListener interface {
	ListeningSocket
	// AcceptConn waits for and returns the next connection to the listener. Same as Accept.
	AcceptConn() (socket StreamConn, err error)
}

type StreamConn interface {
	Socket
	net.Conn
}

type PacketConn interface {
	Socket
	net.PacketConn
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
	SupportsDatagrams() (bool, error)
}

type MuxStream interface {
	StreamConn
	StreamID() int
}
