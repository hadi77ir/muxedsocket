package chaining

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/utils"
)

type LayerType int

const (
	LayerNone       = LayerType(0)
	LayerPacketConn = LayerType(0x01)
	LayerStreamConn = LayerType(0x02)
	// Stream-over-packets
	LayerStreamAdapter = LayerType(0x04)
	// Packets-over-streams
	LayerPacketAdapter    = LayerType(0x08)
	LayerStreamObfuscator = LayerType(0x10)
	LayerPacketObfuscator = LayerType(0x20)
	LayerPacketSolution   = LayerType(0x40)
	LayerStreamSolution   = LayerType(0x80)
	LayerAddrSolution     = LayerType(0x100)
)

// which can be underneath this layer? (top-down)
var compatibilityMatrix = map[LayerType]LayerType{
	LayerNone:             LayerNone,
	LayerPacketConn:       LayerNone,
	LayerStreamConn:       LayerNone,
	LayerAddrSolution:     LayerNone,
	LayerPacketAdapter:    LayerStreamObfuscator | LayerStreamAdapter | LayerStreamConn | LayerAddrSolution | LayerPacketSolution | LayerStreamSolution,
	LayerPacketObfuscator: LayerPacketObfuscator | LayerPacketAdapter | LayerPacketConn,
	LayerStreamObfuscator: LayerStreamObfuscator | LayerStreamAdapter | LayerStreamConn | LayerAddrSolution | LayerPacketSolution | LayerStreamSolution,
	LayerStreamAdapter:    LayerPacketObfuscator | LayerPacketAdapter | LayerPacketConn,
	LayerPacketSolution:   LayerPacketObfuscator | LayerPacketAdapter | LayerPacketConn,
	LayerStreamSolution:   LayerStreamObfuscator | LayerStreamAdapter | LayerStreamConn,
}

const LayersTakingStreamConn = LayerPacketAdapter | LayerStreamObfuscator | LayerStreamSolution
const LayersTakingPacketConn = LayerStreamAdapter | LayerPacketObfuscator | LayerPacketSolution
const LayersTakingString = LayerStreamConn | LayerPacketConn | LayerAddrSolution

// GetCompatibleLayers returns compatible layers IDs ORed together.
func GetCompatibleLayers(code LayerType) LayerType {
	return compatibilityMatrix[code]
}

// which layers implement the mux?
const muxerLayers = LayerStreamSolution | LayerPacketSolution | LayerAddrSolution

// GetMuxerLayers returns IDs for layers that implement multiplexer, ORed together.
func GetMuxerLayers() LayerType {
	return muxerLayers
}

func IsMuxerLayer(layerType LayerType) bool {
	return muxerLayers&layerType == layerType
}

// which layers implement the basic connection?
const transportLayers = LayerPacketConn | LayerStreamConn | LayerAddrSolution

func GetTransportLayers() LayerType {
	return transportLayers
}

func IsTransportLayer(layerType LayerType) bool {
	return transportLayers&layerType == layerType
}

type Layer struct {
	LayerType          LayerType
	ImplementationName string
	Implementation     any
	ParametersIndex    int
}

func IsLayerCompatible(above, below *Layer) bool {
	if below == nil {
		return IsInterfaceCompatible(above.LayerType, 0)
	}
	return IsInterfaceCompatible(above.LayerType, below.LayerType)
}

func IsInterfaceCompatible(above, below LayerType) bool {
	compatibleTypes := compatibilityMatrix[above]
	if compatibleTypes == LayerNone {
		return below == LayerNone
	}
	if compatibleTypes&below == 0x0 {
		return false
	}
	return true
}

func tryFindPart[T any](registry *utils.Registry[T], part string, paramsSectionIndex int, typeId LayerType) (*Layer, bool) {
	impl, found := registry.Get(part)
	if found {
		return &Layer{
			LayerType:          typeId,
			ImplementationName: part,
			ParametersIndex:    paramsSectionIndex,
			Implementation:     impl,
		}, true
	}
	return nil, false
}

func ResolveLayer(creators *muxedsocket.Creators, part string, paramsSectionIndex int, enableMux bool) *Layer {
	if enableMux {
		addrSolution, found := tryFindPart(creators.AddrSolutions(), part, paramsSectionIndex, LayerAddrSolution)
		if found {
			return addrSolution
		}
		streamSolution, found := tryFindPart(creators.StreamSolutions(), part, paramsSectionIndex, LayerStreamSolution)
		if found {
			return streamSolution
		}
		packetSolution, found := tryFindPart(creators.PacketSolutions(), part, paramsSectionIndex, LayerPacketSolution)
		if found {
			return packetSolution
		}
	}
	packetObfuscator, found := tryFindPart(creators.PacketObfuscators(), part, paramsSectionIndex, LayerPacketObfuscator)
	if found {
		return packetObfuscator
	}
	streamObfuscator, found := tryFindPart(creators.StreamObfuscators(), part, paramsSectionIndex, LayerStreamObfuscator)
	if found {
		return streamObfuscator
	}
	packetAdapter, found := tryFindPart(creators.PacketAdapters(), part, paramsSectionIndex, LayerPacketAdapter)
	if found {
		return packetAdapter
	}
	streamAdapter, found := tryFindPart(creators.StreamAdapters(), part, paramsSectionIndex, LayerStreamAdapter)
	if found {
		return streamAdapter
	}
	streamConn, found := tryFindPart(creators.StreamConns(), part, paramsSectionIndex, LayerStreamConn)
	if found {
		return streamConn
	}
	packetConn, found := tryFindPart(creators.PacketConns(), part, paramsSectionIndex, LayerPacketConn)
	if found {
		return packetConn
	}

	return nil
}

func splitParameters(layers []*Layer, parameters utils.Parameters) []utils.Parameters {
	layerParams := make([]utils.Parameters, len(layers))
	for i := 0; i < len(layers); i++ {
		layerParams[i] = parameters.SectionWithCommon(utils.GetIndexedParamsSection(layers[i].ImplementationName, layers[i].ParametersIndex))
	}
	return layerParams
}

func ResolveLayers(creators *muxedsocket.Creators, schemeParts []string, enableMux bool) ([]*Layer, error) {
	layers := make([]*Layer, len(schemeParts))
	occurrenceMap := make(map[string]int)
	for i := len(schemeParts) - 1; i >= 0; i-- {
		part := schemeParts[i]
		resolved := ResolveLayer(creators, part, occurrenceMap[part], enableMux)
		if resolved == nil {
			return nil, muxedsocket.ErrSchemeNotSupported
		}
		//bottom := getLayerAt(layers, i+1)
		//if bottom != nil && !IsLayerCompatible(resolved, bottom) {
		// todo: we may try to fill the gap with some defaults.
		//	return nil, muxedsocket.ErrIncompatibleChainOfLayers
		//}
		occurrenceMap[part] += 1
		layers[i] = resolved
	}
	return layers, nil
}

func getLayerAt(layers []*Layer, at int) *Layer {
	if at < 0 {
		if len(layers)+at >= len(layers) {
			return nil
		}
		return layers[len(layers)+at]
	}
	if at >= len(layers) {
		return nil
	}
	return layers[at]
}
