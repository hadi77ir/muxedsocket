package yamux

import (
	"github.com/hadi77ir/muxedsocket/types"
	Y "github.com/hashicorp/yamux"
	"net"
)

func wrapConn(conn *Y.Session) types.MuxedSocket {
	return &Conn{
		session: conn,
	}
}

func wrapListener(listener net.Listener, config *Y.Config) *Listener {
	return &Listener{listener: listener, config: config}
}

func wrapStream(stream *Y.Stream, session *Y.Session) types.MuxStream {
	return &Stream{
		stream:     stream,
		localAddr:  types.WrapAddr(session.LocalAddr(), int(stream.StreamID())),
		remoteAddr: types.WrapAddr(session.RemoteAddr(), int(stream.StreamID())),
	}
}
