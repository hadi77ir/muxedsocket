package packet

import (
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
	"net"
)

type StandardPacketConnFunc func(network, addr string) (net.PacketConn, error)
type StandardPrimedPacketConnFunc func() (net.PacketConn, error)
type StandardResolveFunc func(network, addr string) (net.Addr, error)
type AfterConnectHookFunc func(conn net.PacketConn, parameters utils.Parameters)
type WrappedHookFunc func(conn net.PacketConn)

type PacketConnImplementation struct {
	dialFunc        StandardPacketConnFunc
	listenFunc      StandardPacketConnFunc
	resolveFunc     StandardResolveFunc
	network         string
	afterDialHook   *utils.Hook[AfterConnectHookFunc]
	afterListenHook *utils.Hook[AfterConnectHookFunc]
}

var _ types.PacketConnImplementation = &PacketConnImplementation{}

func (p *PacketConnImplementation) Server(addr string, parameters utils.Parameters) (types.PacketConnFunc, error) {
	afterDial := WrapAfterConnectHooksFunc(p.afterDialHook, parameters)
	dialer := WrapPacketConnFunc(p.dialFunc, p.network, addr)
	return func() (types.PacketConn, error) {
		conn, err := dialer()
		if err != nil {
			return nil, err
		}
		afterDial(conn)
		return WrapConn(conn, nil, dialer, afterDial), nil
	}, nil
}

func (p *PacketConnImplementation) Client(addr string, parameters utils.Parameters) (types.PacketConnFunc, error) {
	afterDial := WrapAfterConnectHooksFunc(p.afterDialHook, parameters)
	remoteAddr, err := p.resolveFunc(p.network, addr)
	if err != nil {
		return nil, err
	}
	dialer := WrapPacketConnFunc(p.dialFunc, p.network, remoteAddr.String())
	return func() (types.PacketConn, error) {
		conn, err := dialer()
		if err != nil {
			return nil, err
		}
		afterDial(conn)
		return WrapConn(conn, remoteAddr, dialer, afterDial), nil
	}, nil
}

func (n *PacketConnImplementation) AfterListenHook() *utils.Hook[AfterConnectHookFunc] {
	return n.afterListenHook
}

func (n *PacketConnImplementation) AfterDialHook() *utils.Hook[AfterConnectHookFunc] {
	return n.afterDialHook
}

func WrapPacketConnFunc(dialFunc StandardPacketConnFunc, network string, addr string) StandardPrimedPacketConnFunc {
	return func() (net.PacketConn, error) {
		return dialFunc(network, addr)
	}
}

func WrapAfterConnectHooksFunc(hook *utils.Hook[AfterConnectHookFunc], parameters utils.Parameters) WrappedHookFunc {
	return func(conn net.PacketConn) {
		if hook != nil {
			callHookFuncs(hook, conn, parameters)
		}
	}
}

func callHookFuncs(hook *utils.Hook[AfterConnectHookFunc], conn net.PacketConn, parameters utils.Parameters) {
	for _, hEntry := range hook.GetEntries() {
		if hEntry != nil {
			hFunc := hEntry.GetFunc()
			if hFunc != nil {
				hFunc(conn, parameters)
			}
		}
	}
}

func WrapImplementation(network string, dialFunc, listenFunc StandardPacketConnFunc, resolveFunc StandardResolveFunc) *PacketConnImplementation {
	return &PacketConnImplementation{
		network:         network,
		dialFunc:        dialFunc,
		listenFunc:      listenFunc,
		resolveFunc:     resolveFunc,
		afterListenHook: utils.NewHook[AfterConnectHookFunc](),
		afterDialHook:   utils.NewHook[AfterConnectHookFunc](),
	}
}

func WrapConn(conn net.PacketConn, remoteAddr net.Addr, dialFunc StandardPrimedPacketConnFunc, afterDial WrappedHookFunc) types.PacketConn {
	if oobConn, ok := conn.(OOBCapablePacketConn); ok {
		return wrapOOBConn(oobConn, remoteAddr, dialFunc, afterDial)
	}
	return wrapRawConn(conn, remoteAddr, dialFunc, afterDial)
}
