package stream

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
	"net"
	"time"
)

type StandardDialFunc func(network, addr string) (net.Conn, error)
type StandardDialTimeoutFunc func(network, addr string, timeout time.Duration) (net.Conn, error)
type StandardListenFunc func(network, addr string) (net.Listener, error)
type StandardPrimedDialFunc func() (net.Conn, error)
type AfterConnectHookFunc func(conn net.Conn, parameters utils.Parameters)
type WrappedHookFunc func(conn net.Conn)

type NetConnImplementation struct {
	dialFunc        StandardDialTimeoutFunc
	listenFunc      StandardListenFunc
	network         string
	afterDialHook   *utils.Hook[AfterConnectHookFunc]
	afterAcceptHook *utils.Hook[AfterConnectHookFunc]
}

var _ types.StreamConnImplementation = &NetConnImplementation{}

func (n *NetConnImplementation) Server(addr string, parameters utils.Parameters) (types.StreamListenFunc, error) {
	return func() (types.StreamListener, error) {
		listener, err := n.listenFunc(n.network, addr)
		if err != nil {
			return nil, err
		}
		return WrapListener(listener, WrapAfterConnectHooksFunc(n.afterAcceptHook, parameters)), err
	}, nil
}

func (n *NetConnImplementation) Client(addr string, parameters utils.Parameters) (types.StreamDialFunc, error) {
	dialer := WrapDialTimeoutFunc(n.dialFunc, n.network, addr, utils.DurationFromParameters(parameters, muxedsocket.ParamDialTimeout, muxedsocket.DefaultDialTimeout))
	afterDial := WrapAfterConnectHooksFunc(n.afterDialHook, parameters)
	return func() (types.StreamConn, error) {
		conn, err := dialer()
		if err != nil {
			return nil, err
		}
		afterDial(conn)
		return WrapConn(conn, dialer, afterDial), nil
	}, nil
}

func WrapAfterConnectHooksFunc(hook *utils.Hook[AfterConnectHookFunc], parameters utils.Parameters) WrappedHookFunc {
	return func(conn net.Conn) {
		if hook != nil {
			callHookFuncs(hook, conn, parameters)
		}
	}
}

func callHookFuncs(hook *utils.Hook[AfterConnectHookFunc], conn net.Conn, parameters utils.Parameters) {
	for _, hEntry := range hook.GetEntries() {
		if hEntry != nil {
			hFunc := hEntry.GetFunc()
			if hFunc != nil {
				hFunc(conn, parameters)
			}
		}
	}
}

func (n *NetConnImplementation) AfterAcceptHook() *utils.Hook[AfterConnectHookFunc] {
	return n.afterAcceptHook
}

func (n *NetConnImplementation) AfterDialHook() *utils.Hook[AfterConnectHookFunc] {
	return n.afterDialHook
}

func DialTimeoutAdapter(dialFunc StandardDialFunc) StandardDialTimeoutFunc {
	return func(network, addr string, _ time.Duration) (net.Conn, error) {
		return dialFunc(network, addr)
	}
}

func WrapDialTimeoutFunc(dialFunc StandardDialTimeoutFunc, network string, addr string, timeout time.Duration) StandardPrimedDialFunc {
	return func() (net.Conn, error) {
		return dialFunc(network, addr, timeout)
	}
}

func WrapStandardImplementation(network string, dialFunc StandardDialTimeoutFunc, listenFunc StandardListenFunc) *NetConnImplementation {
	return &NetConnImplementation{
		dialFunc:        dialFunc,
		listenFunc:      listenFunc,
		network:         network,
		afterDialHook:   utils.NewHook[AfterConnectHookFunc](),
		afterAcceptHook: utils.NewHook[AfterConnectHookFunc](),
	}
}

type WrappedListener struct {
	listener    net.Listener
	closed      chan struct{}
	afterAccept func(conn net.Conn)
}

func (w *WrappedListener) CloseChan() <-chan struct{} {
	return w.closed
}

func (w *WrappedListener) Close() error {
	select {
	case <-w.closed:
		break
	default:
		close(w.closed)
	}
	return w.Close()
}

func (w *WrappedListener) Accept() (socket types.Socket, err error) {
	return w.AcceptConn()
}

func (w *WrappedListener) Addr() net.Addr {
	return w.listener.Addr()
}

func (w *WrappedListener) AcceptConn() (socket types.StreamConn, err error) {
	conn, err := w.listener.Accept()
	if err != nil {
		return
	}
	if w.afterAccept != nil {
		w.afterAccept(conn)
	}
	socket = WrapConn(conn, nil, nil)
	return
}

func WrapListener(listener net.Listener, afterAccept WrappedHookFunc) types.StreamListener {
	return &WrappedListener{listener: listener, closed: make(chan struct{}, 1), afterAccept: afterAccept}
}

type WrappedConn struct {
	net.Conn
	dialFunc    StandardPrimedDialFunc
	closed      chan struct{}
	afterRedial WrappedHookFunc
}

func (w *WrappedConn) CloseChan() <-chan struct{} {
	return w.closed
}

func (w *WrappedConn) CanRedial() bool {
	return w.dialFunc != nil
}

func (w *WrappedConn) Redial() (types.Socket, error) {
	if w.dialFunc == nil {
		return nil, muxedsocket.ErrRedialNotSupported
	}
	conn, err := w.dialFunc()
	if err != nil {
		return nil, err
	}
	if w.afterRedial != nil {
		w.afterRedial(conn)
	}
	return WrapConn(conn, w.dialFunc, w.afterRedial), nil
}

func (w *WrappedConn) Read(b []byte) (n int, err error) {
	n, err = w.Read(b)
	w.handleError(err)
	return
}

func (w *WrappedConn) Write(b []byte) (n int, err error) {
	n, err = w.Write(b)
	w.handleError(err)
	return
}

func (w *WrappedConn) Close() error {
	select {
	case <-w.closed:
		break
	default:
		close(w.closed)
	}
	return w.Close()
}

func (w *WrappedConn) handleError(err error) {
	if utils.IsConnEOL(err) {
		_ = w.Close()
	}
}

func WrapConn(conn net.Conn, dialFunc StandardPrimedDialFunc, afterRedial WrappedHookFunc) types.StreamConn {
	return &WrappedConn{Conn: conn, dialFunc: dialFunc, closed: make(chan struct{}, 1), afterRedial: afterRedial}
}

func WrapDialer(dialFunc StandardPrimedDialFunc) types.StreamDialFunc {
	return func() (types.StreamConn, error) {
		conn, err := dialFunc()
		if err != nil {
			return nil, err
		}
		return WrapConn(conn, dialFunc, nil), nil
	}
}
