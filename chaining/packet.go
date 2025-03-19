package chaining

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
)

type packetLayersChainer struct {
	backend *genericChainer
}

func (m *packetLayersChainer) ConstructDialFunc(addr string, parameters utils.Parameters) (types.PacketConnFunc, error) {
	result, commonParams, err := m.backend.ConstructDialFunc(addr, parameters)
	if err != nil {
		return nil, err
	}
	return GetPacketDialFunc(result, m.backend.Defaults, commonParams)
}

func (m *packetLayersChainer) ConstructListenFunc(addr string, parameters utils.Parameters) (types.PacketConnFunc, error) {
	result, commonParams, err := m.backend.ConstructListenFunc(addr, parameters)
	if err != nil {
		return nil, err
	}
	return GetPacketListenFunc(result, m.backend.Defaults, commonParams)
}

var _ PacketLayersChainer = &packetLayersChainer{}

func CreatePacketLayersChainer(creators *muxedsocket.Creators, defaults *muxedsocket.DefaultLayers, schemeParts []string) (PacketLayersChainer, error) {
	backend, err := createGenericChainer(creators, defaults, schemeParts)
	if err != nil {
		return nil, err
	}
	return &packetLayersChainer{backend: backend}, nil
}
