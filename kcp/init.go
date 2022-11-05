package kcp

import (
	"github.com/hadi77ir/muxedsocket"
	K "github.com/xtaci/kcp-go/v5"
	"net"
)

const (
	parityShards = 3
	dataShards   = 10
)

// Note that this KCP implementation doesn't have any encryption in transport layer and uses NoneBlockCipher as its cipher.
// TODO: Make changes to the signature to support providing block ciphers, data shards and parity shards as parameters.

func ClientAdapter(conn net.PacketConn, remoteAddr string) (net.Conn, error) {
	crypt, _ := K.NewNoneBlockCrypt([]byte{0})
	return K.NewConn(remoteAddr, crypt, dataShards, parityShards, conn)
}

func ServerAdapter(conn net.PacketConn) (net.Listener, error) {
	crypt, _ := K.NewNoneBlockCrypt([]byte{0})
	return K.ServeConn(crypt, dataShards, parityShards, conn)
}

func init() {
	muxedsocket.GlobalCreators().ClientStreamAdapters().Register("kcp", ClientAdapter)
	muxedsocket.GlobalCreators().ServerStreamAdapters().Register("kcp", ServerAdapter)
}
