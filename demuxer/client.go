package demuxer

import (
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
	"sync"
)

func DemuxDialer(dialFunc types.MuxDialFunc, parameters utils.Parameters) types.StreamDialFunc {
	demuxer := CreateClientDemuxer(dialFunc, parameters)
	return demuxer.Dial
}

func CreateClientDemuxer(dialFunc types.MuxDialFunc, parameters utils.Parameters) *ClientDemuxer {
	return &ClientDemuxer{
		dialFunc:             dialFunc,
		streamsPerConnection: utils.IntegerFromParameters(parameters, "streamsperconn", 1),
	}
}

type ClientDemuxer struct {
	streamsPerConnection int
	dialFunc             types.MuxDialFunc
	connection           types.MuxedSocket
	connectionMutex      *sync.Mutex
	currentIteration     int
}

func (d *ClientDemuxer) Dial() (types.StreamConn, error) {
	d.connectionMutex.Lock()
	defer func() {
		d.currentIteration++
		d.connectionMutex.Unlock()
	}()
	d.currentIteration = d.currentIteration % d.streamsPerConnection
	err := d.reconnectIfNeeded()
	if err != nil {
		return nil, err
	}
	return d.connection.OpenStream()
}
func (d *ClientDemuxer) reconnectIfNeeded() error {
	reconnect := false
	if d.currentIteration == 0 {
		reconnect = true
	}
	if d.connection != nil {
		select {
		case <-d.connection.CloseChan():
			reconnect = true
		default:
		}
	}
	if reconnect {
		conn, err := d.dialFunc()
		if err != nil {
			return err
		}
		d.connection = conn
	}
	return nil
}
