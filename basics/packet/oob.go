package packet

import (
	"github.com/hadi77ir/muxedsocket"
	"net"
	"syscall"
)

// OOBCapablePacketConn is a connection that allows the reading of ECN bits from the IP header.
// If the PacketConn passed to Dial or Listen satisfies this interface, quic-go will use it.
// In this case, ReadMsgUDP() will be used instead of ReadFrom() to read packets.
// code from https://github.com/lucas-clemente/quic-go/blob/d2512193dac18f24b31bb3f1f50ac60adb448ee4/sys_conn.go
type OOBCapablePacketConn interface {
	net.PacketConn
	SyscallConn() (syscall.RawConn, error)
	ReadMsgUDP(b, oob []byte) (n, oobn, flags int, addr *net.UDPAddr, err error)
	WriteMsgUDP(b, oob []byte, addr *net.UDPAddr) (n, oobn int, err error)
}

type oobConnWrapper struct {
	*basicWrapper
}

func (c *oobConnWrapper) SyscallConn() (syscall.RawConn, error) {
	if oobConn, ok := c.PacketConn.(OOBCapablePacketConn); ok {
		return oobConn.SyscallConn()
	}
	return nil, muxedsocket.ErrOpNotSupported
}

func (c *oobConnWrapper) ReadMsgUDP(b, oob []byte) (n, oobn, flags int, addr *net.UDPAddr, err error) {
	if oobConn, ok := c.PacketConn.(OOBCapablePacketConn); ok {
		return oobConn.ReadMsgUDP(b, oob)
	}
	return 0, 0, 0, nil, muxedsocket.ErrOpNotSupported
}

func (c *oobConnWrapper) WriteMsgUDP(b, oob []byte, addr *net.UDPAddr) (n, oobn int, err error) {
	if oobConn, ok := c.PacketConn.(OOBCapablePacketConn); ok {
		return oobConn.WriteMsgUDP(b, oob, addr)
	}
	return 0, 0, muxedsocket.ErrOpNotSupported
}

var _ OOBCapablePacketConn = &oobConnWrapper{}

func wrapOOBConn(conn OOBCapablePacketConn, remoteAddr net.Addr, dialFunc StandardPrimedPacketConnFunc, afterDial WrappedHookFunc) *oobConnWrapper {
	return &oobConnWrapper{&basicWrapper{closed: make(chan struct{}, 1), PacketConn: conn, remoteAddr: remoteAddr, dialFunc: dialFunc, afterDial: afterDial}}
}
