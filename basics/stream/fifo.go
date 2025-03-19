package stream

import (
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
	"io"
	"net"
	"sync"
	"time"
)

type ReadFifoDialFunc func() (io.ReadCloser, error)
type WriteFifoDialFunc func() (io.WriteCloser, error)

// DoubleFifoConn contains two FIFOs as input and output.
type DoubleFifoConn struct {
	readFifo  io.ReadCloser
	writeFifo io.WriteCloser
	closed    chan struct{}

	readMutex  sync.Mutex
	writeMutex sync.Mutex
}

func (c *DoubleFifoConn) CloseChan() <-chan struct{} {
	return c.closed
}
func (c *DoubleFifoConn) Close() error {
	select {
	case <-c.closed:
		break
	default:
		close(c.closed)
	}
	rErr := c.readFifo.Close()
	wErr := c.writeFifo.Close()
	if rErr != nil {
		return rErr
	}
	return wErr
}

func (c *DoubleFifoConn) LocalAddr() net.Addr {
	return types.EmptyAddr("fifo:local")
}

func (c *DoubleFifoConn) RemoteAddr() net.Addr {
	return types.EmptyAddr("fifo:remote")
}

func (c *DoubleFifoConn) Read(b []byte) (n int, err error) {
	c.readMutex.Lock()
	defer c.readMutex.Unlock()
	n, err = c.readFifo.Read(b)
	c.handleError(err)
	return
}

func (c *DoubleFifoConn) Write(b []byte) (n int, err error) {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	n, err = c.writeFifo.Write(b)
	c.handleError(err)
	return
}

func (c *DoubleFifoConn) SetDeadline(t time.Time) error {
	return muxedsocket.ErrOpNotSupported
}

func (c *DoubleFifoConn) SetReadDeadline(t time.Time) error {
	return muxedsocket.ErrOpNotSupported
}

func (c *DoubleFifoConn) SetWriteDeadline(t time.Time) error {
	return muxedsocket.ErrOpNotSupported
}

func (c *DoubleFifoConn) handleError(err error) {
	if utils.IsConnEOL(err) {
		_ = c.Close()
	}
}

func DialDoubleFifo(readDialer ReadFifoDialFunc, writeDialer WriteFifoDialFunc) (*DoubleFifoConn, error) {
	readFifo, err := readDialer()
	if err != nil {
		return nil, err
	}
	writeFifo, err := writeDialer()
	if err != nil {
		_ = readFifo.Close()
		return nil, err
	}
	return WrapFifoConn(readFifo, writeFifo), nil
}

func WrapFifoConn(readFifo io.ReadCloser, writeFifo io.WriteCloser) *DoubleFifoConn {
	return &DoubleFifoConn{
		readFifo:  readFifo,
		writeFifo: writeFifo,
		closed:    make(chan struct{}, 1),
	}
}

var _ net.Conn = &DoubleFifoConn{}

type DoubleFifoListener struct {
	closed          chan struct{}
	readFifoDialer  ReadFifoDialFunc
	writeFifoDialer WriteFifoDialFunc
	availableToDial <-chan struct{}
}

func (c *DoubleFifoListener) Accept() (net.Conn, error) {
	select {
	case <-c.availableToDial:
		break
	case <-c.closed:
		return nil, net.ErrClosed
	}
	socket, err := DialDoubleFifo(c.readFifoDialer, c.writeFifoDialer)
	if err != nil {
		return nil, err
	}
	c.availableToDial = socket.CloseChan()
	return socket, nil
}

func (c *DoubleFifoListener) Addr() net.Addr {
	return types.EmptyAddr("fifo:listener")
}

func (c *DoubleFifoListener) CloseChan() <-chan struct{} {
	return c.closed
}

func (c *DoubleFifoListener) Close() error {
	select {
	case <-c.closed:
		break
	default:
		close(c.closed)
	}
	return nil
}

func NewFifoListener(readFifoDialer ReadFifoDialFunc, writeFifoDialer WriteFifoDialFunc) *DoubleFifoListener {
	closedChan := make(chan struct{}, 1)
	close(closedChan)
	return &DoubleFifoListener{
		availableToDial: closedChan,
		closed:          make(chan struct{}, 1),
		readFifoDialer:  readFifoDialer,
		writeFifoDialer: writeFifoDialer,
	}
}

var _ net.Listener = &DoubleFifoListener{}
