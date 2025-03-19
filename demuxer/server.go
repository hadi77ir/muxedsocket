package demuxer

import (
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
	"net"
	"sync"
	"sync/atomic"
)

func DemuxListener(listenFunc types.MuxListenFunc, parameters utils.Parameters) types.StreamListenFunc {
	demuxer := CreateServerDemuxer(listenFunc, parameters)
	return demuxer.Listen
}

func CreateServerDemuxer(listenFunc types.MuxListenFunc, parameters utils.Parameters) *ServerDemuxer {
	return &ServerDemuxer{
		listenFunc:          listenFunc,
		streamAcceptBacklog: utils.IntegerFromParameters(parameters, "backlog", 1000),
	}
}

type ServerDemuxer struct {
	listenFunc          types.MuxListenFunc
	streamAcceptBacklog int
}

func (d *ServerDemuxer) Listen() (types.StreamListener, error) {
	listener, err := d.listenFunc()
	if err != nil {
		return nil, err
	}
	listenerDemuxer := &ListenerDemuxer{
		listener: listener,
	}
	listenerDemuxer.StartWorker(d.streamAcceptBacklog)
	return listenerDemuxer, nil
}

type ListenerDemuxer struct {
	listener types.MuxedListener
	wg       *sync.WaitGroup
	backlog  chan types.StreamConn
	done     chan struct{}
	init     atomic.Bool
	running  atomic.Bool
}

func (l *ListenerDemuxer) CloseChan() <-chan struct{} {
	return l.listener.CloseChan()
}

func (l *ListenerDemuxer) Close() error {
	l.signalClose()
	return l.listener.Close()
}

func (l *ListenerDemuxer) Accept() (socket types.Socket, err error) {
	return l.AcceptConn()
}

func (l *ListenerDemuxer) Addr() net.Addr {
	return l.listener.Addr()
}

func (l *ListenerDemuxer) AcceptConn() (socket types.StreamConn, err error) {
	select {
	case c := <-l.backlog:
		return c, nil
	case <-l.listener.CloseChan():
		return nil, net.ErrClosed
	case <-l.done:
		return nil, net.ErrClosed
	}
}

func (l *ListenerDemuxer) StartWorker(backlog int) {
	l.tryInit()
	oldState := l.running.Swap(true)
	if oldState == false {
		l.backlog = make(chan types.StreamConn, backlog)
		l.wg.Add(1)
		go l.acceptMuxedWorker()
	}
}

func (l *ListenerDemuxer) tryInit() {
	oldState := l.init.Swap(true)
	if oldState == false {
		l.wg = &sync.WaitGroup{}
		l.done = make(chan struct{}, 1)
	}
}

func (l *ListenerDemuxer) acceptMuxedWorker() {
	defer l.wg.Done()
	for {
		accepted, err := l.listener.Accept()
		if err != nil {
			l.signalClose()
			return
		}
		l.wg.Add(1)
		go l.acceptStreamsWorker(accepted.(types.MuxedSocket))
	}
}

func (l *ListenerDemuxer) acceptStreamsWorker(sock types.MuxedSocket) {
	defer l.wg.Done()
	for {
		select {
		case <-l.done:
			return
		case <-l.listener.CloseChan():
			return
		default:
		}
		stream, err := sock.AcceptStream()
		if err != nil {
			select {
			case <-sock.CloseChan():
				return
			default:
			}
		}
		l.backlog <- stream
	}
}

func (l *ListenerDemuxer) signalClose() {
	select {
	case <-l.done:
	case <-l.listener.CloseChan():
	default:
		close(l.done)
	}
}
