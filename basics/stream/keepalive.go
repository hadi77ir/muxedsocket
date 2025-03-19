package stream

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/utils"
	"net"
	"time"
)

// KeepAliveCapableConn defines the portion of TCPConn that is related to keep-alive functionality.
// Taken from Go Standard library (net/tcpsock.go)
type KeepAliveCapableConn interface {
	// SetKeepAlive sets whether the operating system should send
	// keep-alive messages on the connection.
	SetKeepAlive(keepalive bool) error

	// SetKeepAlivePeriod sets period between keep-alives.
	SetKeepAlivePeriod(d time.Duration) error
}

func AddKeepAliveHook(impl *NetConnImplementation) {
	impl.AfterDialHook().Add("keepalive", keepAliveHook)
	impl.AfterAcceptHook().Add("keepalive", keepAliveHook)
}

func keepAliveHook(conn net.Conn, parameters utils.Parameters) {
	if kacConn, ok := conn.(KeepAliveCapableConn); ok {
		enableKeepAlive := true
		keepAlivePeriod := muxedsocket.DefaultKeepAlive
		if keepAlive, found := parameters.Get(muxedsocket.ParamKeepAlive); found {
			if utils.StrIsFalse(keepAlive) {
				enableKeepAlive = false
			}
			var err error
			keepAlivePeriod, err = time.ParseDuration(keepAlive)
			if err != nil {
				keepAlivePeriod = muxedsocket.DefaultKeepAlive
			}
		}
		_ = kacConn.SetKeepAlive(enableKeepAlive)
		_ = kacConn.SetKeepAlivePeriod(keepAlivePeriod)
	}
}
