package chaining

import (
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
)

type LayersChainer[TDialFunc any, TListenFunc any] interface {
	ConstructDialFunc(addr string, parameters utils.Parameters) (TDialFunc, error)
	ConstructListenFunc(addr string, parameters utils.Parameters) (TListenFunc, error)
}

type MuxLayersChainer LayersChainer[types.MuxDialFunc, types.MuxListenFunc]
type StreamLayersChainer LayersChainer[types.StreamDialFunc, types.StreamListenFunc]
type PacketLayersChainer LayersChainer[types.PacketConnFunc, types.PacketConnFunc]
