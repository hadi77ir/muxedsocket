package chaining

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
)

type muxLayersChainer struct {
	backend *genericChainer
}

func (m *muxLayersChainer) ConstructDialFunc(addr string, parameters utils.Parameters) (types.MuxDialFunc, error) {
	result, commonParams, err := m.backend.ConstructDialFunc(addr, parameters)
	if err != nil {
		return nil, err
	}
	return GetMuxDialFunc(result, m.backend.Defaults, commonParams)
}

func (m *muxLayersChainer) ConstructListenFunc(addr string, parameters utils.Parameters) (types.MuxListenFunc, error) {
	result, commonParams, err := m.backend.ConstructListenFunc(addr, parameters)
	if err != nil {
		return nil, err
	}
	return GetMuxListenFunc(result, m.backend.Defaults, commonParams)
}

var _ MuxLayersChainer = &muxLayersChainer{}

func CreateMuxLayersChainer(creators *muxedsocket.Creators, defaults *muxedsocket.DefaultLayers, schemeParts []string) (MuxLayersChainer, error) {
	backend, err := createGenericChainer(creators, defaults, schemeParts)
	if err != nil {
		return nil, err
	}
	return &muxLayersChainer{backend: backend}, nil
}
