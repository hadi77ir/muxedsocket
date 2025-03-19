package types

import (
	"github.com/hadi77ir/muxedsocket/utils"
)

type ConcreteStreamListenFunc func(addr string, parameters utils.Parameters) (MuxListenFunc, error)
type ConcreteStreamDialFunc func(addr string, parameters utils.Parameters) (MuxDialFunc, error)

type PacketConnFunc func() (PacketConn, error)
type StreamDialFunc func() (StreamConn, error)
type StreamListenFunc func() (StreamListener, error)
type MuxListenFunc func() (MuxedListener, error)
type MuxDialFunc func() (MuxedSocket, error)
type MuxStreamConnectFunc func() (MuxStream, error)

type PacketConnImplementation interface {
	Server(addr string, parameters utils.Parameters) (PacketConnFunc, error)
	Client(addr string, parameters utils.Parameters) (PacketConnFunc, error)
}

type StreamConnImplementation interface {
	Server(addr string, parameters utils.Parameters) (StreamListenFunc, error)
	Client(addr string, parameters utils.Parameters) (StreamDialFunc, error)
}

type AddrSolutionImplementation interface {
	Server(addr string, parameters utils.Parameters) (MuxListenFunc, error)
	Client(addr string, parameters utils.Parameters) (MuxDialFunc, error)
}

type StreamAdapterImplementation interface {
	Server(conn PacketConnFunc, parameters utils.Parameters) (StreamListenFunc, error)
	Client(conn PacketConnFunc, parameters utils.Parameters) (StreamDialFunc, error)
}

type PacketAdapterImplementation interface {
	Server(listener StreamListenFunc, parameters utils.Parameters) (PacketConnFunc, error)
	Client(dialFunc StreamDialFunc, parameters utils.Parameters) (PacketConnFunc, error)
	SupportsParallel() bool
}

type StreamObfuscatorImplementation interface {
	Server(conn StreamListenFunc, parameters utils.Parameters) (StreamListenFunc, error)
	Client(conn StreamDialFunc, parameters utils.Parameters) (StreamDialFunc, error)
}

type PacketObfuscatorImplementation interface {
	Server(conn PacketConnFunc, parameters utils.Parameters) (PacketConnFunc, error)
	Client(conn PacketConnFunc, parameters utils.Parameters) (PacketConnFunc, error)
}

type PacketSolutionImplementation interface {
	Server(conn PacketConnFunc, parameters utils.Parameters) (MuxListenFunc, error)
	Client(conn PacketConnFunc, parameters utils.Parameters) (MuxDialFunc, error)
}

type StreamSolutionImplementation interface {
	Server(conn StreamListenFunc, parameters utils.Parameters) (MuxListenFunc, error)
	Client(conn StreamDialFunc, parameters utils.Parameters) (MuxDialFunc, error)
}

type HasParametersHint interface {
	ClientParametersHint() []utils.ParameterHint
	ServerParametersHint() []utils.ParameterHint
}
