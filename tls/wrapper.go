package tls

import (
	"github.com/hadi77ir/muxedsocket/basics/stream"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
	"net"
)

type StandardObfuscatorFunc func(net.Conn, any) net.Conn
type StandardParseConfigFunc func(parameters utils.Parameters, isClient bool) (any, error)

type TLSImplementation struct {
	clientFunc       StandardObfuscatorFunc
	serverFunc       StandardObfuscatorFunc
	configParserFunc StandardParseConfigFunc
}

func (p *TLSImplementation) Server(conn types.StreamListenFunc, parameters utils.Parameters) (types.StreamListenFunc, error) {
	params, err := p.configParserFunc(parameters, false)
	if err != nil {
		return nil, err
	}
	return func() (types.StreamListener, error) {
		listener, err := conn()
		if err != nil {
			return nil, err
		}
		return WrapListener(listener, p.serverFunc, params)
	}, nil
}

type ListenerWrapper struct {
	listener   types.StreamListener
	serverFunc StandardObfuscatorFunc
	params     any
}

func (w *ListenerWrapper) CloseChan() <-chan struct{} {
	return w.CloseChan()
}

func (w *ListenerWrapper) Close() error {
	return w.Close()
}

func (w *ListenerWrapper) Accept() (socket types.Socket, err error) {
	return w.AcceptConn()
}

func (w *ListenerWrapper) Addr() net.Addr {
	return w.listener.Addr()
}

func (w *ListenerWrapper) AcceptConn() (socket types.StreamConn, err error) {
	conn, err := w.listener.AcceptConn()
	if err != nil {
		return nil, err
	}
	return stream.WrapConn(w.serverFunc(conn, w.params), nil, nil), nil
}

func WrapListener(listener types.StreamListener, serverFunc StandardObfuscatorFunc, params any) (types.StreamListener, error) {
	return &ListenerWrapper{listener: listener, serverFunc: serverFunc, params: params}, nil
}

func (p *TLSImplementation) Client(conn types.StreamDialFunc, parameters utils.Parameters) (types.StreamDialFunc, error) {
	params, err := p.configParserFunc(parameters, true)
	if err != nil {
		return nil, err
	}
	return wrapDialer(conn, p.clientFunc, params), nil
}

func createTLSDialer(conn types.StreamDialFunc, obfuscatorFunc StandardObfuscatorFunc, parameters any) stream.StandardPrimedDialFunc {
	return func() (net.Conn, error) {
		c, err := conn()
		if err != nil {
			return nil, err
		}
		return obfuscatorFunc(c, parameters), nil
	}
}
func wrapDialer(conn types.StreamDialFunc, obfuscatorFunc StandardObfuscatorFunc, parameters any) types.StreamDialFunc {
	dialer := createTLSDialer(conn, obfuscatorFunc, parameters)
	return func() (types.StreamConn, error) {
		conn, err := dialer()
		if err != nil {
			return nil, err
		}
		return stream.WrapConn(conn, dialer, nil), nil
	}
}

var _ types.StreamObfuscatorImplementation = &TLSImplementation{}

func WrapImplementation(client, server StandardObfuscatorFunc, configFunc StandardParseConfigFunc) *TLSImplementation {
	return &TLSImplementation{clientFunc: client, serverFunc: server, configParserFunc: configFunc}
}

func NewTLSImplementation() types.StreamObfuscatorImplementation {
	return WrapImplementation(ClientTLS, ServerTLS, ParseTLS)
}
