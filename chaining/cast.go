package chaining

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/demuxer"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
)

func GetPacketDialFunc(input any, defaults *muxedsocket.DefaultLayers, commonParameters utils.Parameters) (types.PacketConnFunc, error) {
	if packetConn, ok := input.(types.PacketConnFunc); ok {
		return packetConn, nil
	}
	if streamConn, ok := input.(types.StreamDialFunc); ok {
		// it is a streaming conn. add packets-over-streams implementation.
		// todo: check if stream supports packet transmission. (useful for eNet)
		return defaults.PacketAdapter.Client(streamConn, commonParameters)
	}
	if muxDialer, ok := input.(types.MuxDialFunc); ok {
		// it is a mux dialer. first demux, then packet adapter.
		// todo: check if mux supports packet transmission. (useful for quic)
		return defaults.PacketAdapter.Client(demuxer.DemuxDialer(muxDialer, commonParameters), commonParameters)
	}
	if addr, ok := input.(string); ok {
		// call default stream transport
		return defaults.PacketConn.Client(addr, commonParameters)
	}
	return nil, muxedsocket.ErrInvalidChainingResult
}

func GetStreamDialFunc(input any, defaults *muxedsocket.DefaultLayers, commonParameters utils.Parameters) (types.StreamDialFunc, error) {
	if streamDialer, ok := input.(types.StreamDialFunc); ok {
		return streamDialer, nil
	}
	if muxDialer, ok := input.(types.MuxDialFunc); ok {
		// it is a mux dialer. so demux.
		return demuxer.DemuxDialer(muxDialer, commonParameters), nil
	}
	if packetConn, ok := input.(types.PacketConnFunc); ok {
		// it is a packetConn. add stream-over-packets implementation.
		return defaults.StreamAdapter.Client(packetConn, commonParameters)
	}
	if addr, ok := input.(string); ok {
		// call default stream transport
		return defaults.StreamConn.Client(addr, commonParameters)
	}
	return nil, muxedsocket.ErrInvalidChainingResult
}

func GetMuxDialFunc(input any, defaults *muxedsocket.DefaultLayers, commonParameters utils.Parameters) (types.MuxDialFunc, error) {
	if muxDialer, ok := input.(types.MuxDialFunc); ok {
		// it is a mux dialer. so demux.
		return muxDialer, nil
	}
	if streamDialer, ok := input.(types.StreamDialFunc); ok {
		return defaults.StreamSolution.Client(streamDialer, commonParameters)
	}
	if packetConn, ok := input.(types.PacketConnFunc); ok {
		// it is a packetConn. add stream-over-packets implementation.
		return defaults.PacketSolution.Client(packetConn, commonParameters)
	}
	return nil, muxedsocket.ErrInvalidChainingResult
}

func GetMuxListenFunc(input any, defaults *muxedsocket.DefaultLayers, commonParameters utils.Parameters) (types.MuxListenFunc, error) {
	if muxListen, ok := input.(types.MuxListenFunc); ok {
		return muxListen, nil
	}
	if streamListen, ok := input.(types.StreamListenFunc); ok {
		return defaults.StreamSolution.Server(streamListen, commonParameters)
	}
	if packetConn, ok := input.(types.PacketConnFunc); ok {
		return defaults.PacketSolution.Server(packetConn, commonParameters)
	}
	return nil, muxedsocket.ErrInvalidChainingResult
}

func GetStreamListenFunc(input any, defaults *muxedsocket.DefaultLayers, commonParameters utils.Parameters) (types.StreamListenFunc, error) {
	if streamListen, ok := input.(types.StreamListenFunc); ok {
		return streamListen, nil
	}
	if muxListen, ok := input.(types.MuxListenFunc); ok {
		return demuxer.DemuxListener(muxListen, commonParameters), nil
	}
	if packetConn, ok := input.(types.PacketConnFunc); ok {
		return defaults.StreamAdapter.Server(packetConn, commonParameters)
	}
	return nil, muxedsocket.ErrInvalidChainingResult
}

func GetPacketListenFunc(input any, defaults *muxedsocket.DefaultLayers, commonParameters utils.Parameters) (types.PacketConnFunc, error) {
	if packetConn, ok := input.(types.PacketConnFunc); ok {
		return packetConn, nil
	}
	if muxListen, ok := input.(types.MuxListenFunc); ok {
		// todo: check if stream supports packet transmission. (useful for eNet)
		return defaults.PacketAdapter.Server(demuxer.DemuxListener(muxListen, commonParameters), commonParameters)
	}
	if streamListen, ok := input.(types.StreamListenFunc); ok {
		// todo: check if stream supports packet transmission. (useful for eNet)
		return defaults.PacketAdapter.Server(streamListen, commonParameters)
	}
	return nil, muxedsocket.ErrInvalidChainingResult
}
