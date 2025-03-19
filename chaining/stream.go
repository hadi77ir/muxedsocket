package chaining

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
)

type streamLayersChainer struct {
	backend *genericChainer
}

func (m *streamLayersChainer) ConstructDialFunc(addr string, parameters utils.Parameters) (types.StreamDialFunc, error) {
	result, commonParams, err := m.backend.ConstructDialFunc(addr, parameters)
	if err != nil {
		return nil, err
	}
	return GetStreamDialFunc(result, m.backend.Defaults, commonParams)
}

func (m *streamLayersChainer) ConstructListenFunc(addr string, parameters utils.Parameters) (types.StreamListenFunc, error) {
	result, commonParams, err := m.backend.ConstructListenFunc(addr, parameters)
	if err != nil {
		return nil, err
	}
	return GetStreamListenFunc(result, m.backend.Defaults, commonParams)
}

var _ StreamLayersChainer = &streamLayersChainer{}

func CreateStreamLayersChainer(creators *muxedsocket.Creators, defaults *muxedsocket.DefaultLayers, schemeParts []string) (StreamLayersChainer, error) {
	backend, err := createGenericChainer(creators, defaults, schemeParts)
	if err != nil {
		return nil, err
	}
	return &streamLayersChainer{backend: backend}, nil
}
