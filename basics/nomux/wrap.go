package nomux

import (
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
)

type Implementation struct {
	// nothing.
}

func (i Implementation) Server(conn types.StreamListenFunc, parameters utils.Parameters) (types.MuxListenFunc, error) {
	return func() (types.MuxedListener, error) {
		return WrapServer(conn)
	}, nil
}

func (i Implementation) Client(conn types.StreamDialFunc, parameters utils.Parameters) (types.MuxDialFunc, error) {
	return func() (types.MuxedSocket, error) {
		return WrapDialedConn(conn)
	}, nil
}

var _ types.StreamSolutionImplementation = &Implementation{}
