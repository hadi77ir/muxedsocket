package chaining

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
)

type genericChainer struct {
	Defaults *muxedsocket.DefaultLayers
	Layers   []*Layer
}

func (m *genericChainer) ConstructDialFunc(addr string, parameters utils.Parameters) (any, utils.Parameters, error) {
	return m.chainLayers(addr, parameters, m.applyClientLayerOn)
}

func (m *genericChainer) ConstructListenFunc(addr string, parameters utils.Parameters) (any, utils.Parameters, error) {
	return m.chainLayers(addr, parameters, m.applyServerLayerOn)
}

type applierFn func(layer *Layer, input any, layerParams utils.Parameters, commonParameters utils.Parameters) (any, error)

func (m *genericChainer) applyClientLayerOn(layer *Layer, input any, layerParams utils.Parameters, commonParameters utils.Parameters) (any, error) {
	// if layer takes a stream
	if layer.LayerType&LayersTakingStreamConn != 0 {
		streamDialer, err := GetStreamDialFunc(input, m.Defaults, commonParameters)
		if err != nil {
			return nil, err
		}
		switch layer.LayerType {
		case LayerPacketAdapter:
			layerImpl := layer.Implementation.(types.PacketAdapterImplementation)
			return layerImpl.Client(streamDialer, layerParams)
		case LayerStreamObfuscator:
			layerImpl := layer.Implementation.(types.StreamObfuscatorImplementation)
			return layerImpl.Client(streamDialer, layerParams)
		case LayerStreamSolution:
			layerImpl := layer.Implementation.(types.StreamSolutionImplementation)
			return layerImpl.Client(streamDialer, layerParams)
		}
	}

	// if layer takes a packet
	if layer.LayerType&LayersTakingPacketConn != 0 {
		packetDialer, err := GetPacketDialFunc(input, m.Defaults, commonParameters)
		if err != nil {
			return nil, err
		}
		switch layer.LayerType {
		case LayerStreamAdapter:
			layerImpl := layer.Implementation.(types.StreamAdapterImplementation)
			return layerImpl.Client(packetDialer, layerParams)
		case LayerPacketObfuscator:
			layerImpl := layer.Implementation.(types.PacketObfuscatorImplementation)
			return layerImpl.Client(packetDialer, layerParams)
		case LayerPacketSolution:
			layerImpl := layer.Implementation.(types.PacketSolutionImplementation)
			return layerImpl.Client(packetDialer, layerParams)
		}
	}

	// if layer takes a string
	addr, inputIsString := input.(string)
	if inputIsString && layer.LayerType&LayersTakingString != 0 {
		switch layer.LayerType {
		case LayerStreamConn:
			layerImpl := layer.Implementation.(types.StreamConnImplementation)
			return layerImpl.Client(addr, layerParams)
		case LayerPacketConn:
			layerImpl := layer.Implementation.(types.PacketConnImplementation)
			return layerImpl.Client(addr, layerParams)
		case LayerAddrSolution:
			layerImpl := layer.Implementation.(types.AddrSolutionImplementation)
			return layerImpl.Client(addr, layerParams)
		}
	}
	return nil, muxedsocket.ErrInvalidChainingResult
}

func (m *genericChainer) applyServerLayerOn(layer *Layer, input any, layerParams utils.Parameters, commonParameters utils.Parameters) (any, error) {
	// if layer takes a stream
	if layer.LayerType&LayersTakingStreamConn != 0 {
		streamDialer, err := GetStreamListenFunc(input, m.Defaults, commonParameters)
		if err != nil {
			return nil, err
		}
		switch layer.LayerType {
		case LayerPacketAdapter:
			layerImpl := layer.Implementation.(types.PacketAdapterImplementation)
			return layerImpl.Server(streamDialer, layerParams)
		case LayerStreamObfuscator:
			layerImpl := layer.Implementation.(types.StreamObfuscatorImplementation)
			return layerImpl.Server(streamDialer, layerParams)
		case LayerStreamSolution:
			layerImpl := layer.Implementation.(types.StreamSolutionImplementation)
			return layerImpl.Server(streamDialer, layerParams)
		}
	}

	// if layer takes a packet
	if layer.LayerType&LayersTakingPacketConn != 0 {
		packetDialer, err := GetPacketListenFunc(input, m.Defaults, commonParameters)
		if err != nil {
			return nil, err
		}
		switch layer.LayerType {
		case LayerStreamAdapter:
			layerImpl := layer.Implementation.(types.StreamAdapterImplementation)
			return layerImpl.Server(packetDialer, layerParams)
		case LayerPacketObfuscator:
			layerImpl := layer.Implementation.(types.PacketObfuscatorImplementation)
			return layerImpl.Server(packetDialer, layerParams)
		case LayerPacketSolution:
			layerImpl := layer.Implementation.(types.PacketSolutionImplementation)
			return layerImpl.Server(packetDialer, layerParams)
		}
	}

	// if layer takes a string
	addr, inputIsString := input.(string)
	if inputIsString && layer.LayerType&LayersTakingString != 0 {
		switch layer.LayerType {
		case LayerStreamConn:
			layerImpl := layer.Implementation.(types.StreamConnImplementation)
			return layerImpl.Server(addr, layerParams)
		case LayerPacketConn:
			layerImpl := layer.Implementation.(types.PacketConnImplementation)
			return layerImpl.Server(addr, layerParams)
		case LayerAddrSolution:
			layerImpl := layer.Implementation.(types.AddrSolutionImplementation)
			return layerImpl.Server(addr, layerParams)
		}
	}
	return nil, muxedsocket.ErrInvalidChainingResult
}

func (m *genericChainer) chainLayers(addr string, parameters utils.Parameters, applyLayerOnInputFn applierFn) (any, utils.Parameters, error) {
	var result any
	var err error
	// dial the transport
	layerParams := splitParameters(m.Layers, parameters)
	commonParams := utils.CommonParametersFromMap(parameters)
	result = addr
	for i := len(m.Layers) - 2; i >= 0; i-- {
		result, err = applyLayerOnInputFn(m.Layers[i], result, layerParams[i], commonParams)
		if err != nil {
			return nil, nil, err
		}
	}
	return result, commonParams, nil
}

func createGenericChainer(creators *muxedsocket.Creators, defaults *muxedsocket.DefaultLayers, schemeParts []string) (*genericChainer, error) {
	layers, err := ResolveLayers(creators, schemeParts, false)
	if err != nil {
		return nil, err
	}
	return &genericChainer{Layers: layers, Defaults: defaults}, nil
}
