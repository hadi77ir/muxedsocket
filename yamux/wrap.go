package yamux

import (
	Y "github.com/hashicorp/yamux"
	"muxedsocket"
	"net"
)

func wrapConn(conn *Y.Session) muxedsocket.MuxedSocket {
	return &Conn{
		session: conn,
	}
}

func wrapListener(listener net.Listener, config *Y.Config) *Listener {
	return &Listener{listener: listener, config: config}
}

func wrapStream(stream *Y.Stream, session *Y.Session) muxedsocket.MuxStream {
	return &Stream{
		stream:     stream,
		localAddr:  muxedsocket.WrapAddr(session.LocalAddr(), int(stream.StreamID())),
		remoteAddr: muxedsocket.WrapAddr(session.RemoteAddr(), int(stream.StreamID())),
	}
}
