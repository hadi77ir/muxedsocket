package http

import (
	"errors"
	"github.com/hadi77ir/muxedsocket"
	"github.com/hadi77ir/muxedsocket/basics/stream"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func wrapHttpServer(listener net.Listener, parameters utils.Parameters) (types.StreamListener, error) {
	_, scheme := GetProtocolFromParameters(parameters)
	// cleartext
	if scheme == "http" {
		return wrapListenerH2C(listener, parameters)
	}
	if scheme == "https" {
		return wrapListenerH2S(listener, parameters)
	}
	return nil, muxedsocket.ErrSchemeNotSupported
}

func wrapListenerH2C(listener net.Listener, parameters utils.Parameters) (types.StreamListener, error) {
	handler, backlogChannel, err := newRequestHandler(parameters, ProtoHTTP, plainAdapter)
	if err != nil {
		return nil, err
	}
	server := &http.Server{
		Handler: h2c.NewHandler(handler, &http2.Server{}),
	}
	return &httpListener{
		server:         server,
		backlogChannel: backlogChannel,
		closed:         make(chan struct{}, 1),
	}, nil
}

func wrapListenerH2S(listener net.Listener, parameters utils.Parameters) (types.StreamListener, error) {

}

type requestMatcherFunc func(request *http.Request) bool
type headerMatcherFunc func(header http.Header) bool

func dummyHeaderMatcher(_ http.Header) bool {
	return true
}

func getRequestMatcher(method string, matchingUrl *url.URL, headerMatcher headerMatcherFunc) (requestMatcherFunc, error) {
	if method == "CONNECT" {
		connectUrl, err := createConnectURLFromParameters(matchingUrl.Scheme, matchingUrl.Host)
		if err != nil {
			return nil, err
		}
		return func(request *http.Request) bool {
			return request.Method == http.MethodConnect && request.URL.Host == connectUrl.Host && headerMatcher(request.Header)
		}, nil
	}
	return func(request *http.Request) bool {
		return request.Method == method &&
			request.Host == matchingUrl.Host && request.URL.Path == matchingUrl.Path &&
			matchGetParameters(request.URL.Query(), matchingUrl.Query()) && headerMatcher(request.Header)
	}, nil
}

func matchGetParameters(haystack url.Values, needles url.Values) bool {
	for key, value := range needles {
		if hValues, ok := haystack[key]; ok && hValues != nil && len(hValues) > 0 {
			matches := false
			for _, hValue := range hValues {
				if value[0] == hValue {
					matches = true
					break
				}
			}
			if !matches {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
func plainAdapter(writer http.ResponseWriter, request *http.Request, backlogChan chan net.Conn) {
	var wrappedWriter io.WriteCloser = WriterToWriteCloserAdapter(writer.Write)
	if flusher, ok := writer.(http.Flusher); ok {
		wrappedWriter = WrapWriterAsFlusherWriter(writer, flusher)
	}
	accepted := stream.WrapFifoConn(request.Body, wrappedWriter)
	backlogChan <- accepted
	select {
	case <-accepted.CloseChan():
	case <-request.Context().Done():
	}
}

// requestAdapterFunc transforms given Request and ResponseWriter to a net.Conn and pushes to channel, then waits for
// connection being closed.
type requestAdapterFunc func(w http.ResponseWriter, r *http.Request, c chan net.Conn)

func newRequestHandler(parameters utils.Parameters, scheme string, adapter requestAdapterFunc) (http.Handler, <-chan net.Conn, error) {
	handledUrl := createURLFromParameters(parameters, scheme)
	method := strings.ToUpper(utils.StringFromParameters(parameters, "method", "GET"))
	// todo: extra headers matcher
	// todo: authorization support
	requestMatcher, err := getRequestMatcher(method, handledUrl, dummyHeaderMatcher)
	if err != nil {
		return nil, nil, err
	}
	backlogSize := utils.IntegerFromParameters(parameters, "backlog", 1000)
	if backlogSize < 1 {
		return nil, nil, ErrInvalidBacklogSize
	}
	backlogChan := make(chan net.Conn, backlogSize)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if requestMatcher(request) {
			adapter(writer, request, backlogChan)
			return
		}
		http.NotFound(writer, request)
	}), backlogChan, nil
}

var ErrInvalidBacklogSize = errors.New("backlog size has to be >= 1")

// Workaround for "http.ResponseWriter" not supporting "io.Closer"
type WriterToWriteCloserAdapter func([]byte) (int, error)

func (w WriterToWriteCloserAdapter) Write(p []byte) (int, error) {
	return w(p)
}
func (w WriterToWriteCloserAdapter) Close() error {
	return nil
}

func WrapWriterAsFlusherWriter(writer io.Writer, flusher http.Flusher) io.WriteCloser {
	if flusher != nil {
		return WriterToWriteCloserAdapter(func(p []byte) (int, error) {
			n, err := writer.Write(p)
			flusher.Flush()
			return n, err
		})
	}
	return WriterToWriteCloserAdapter(writer.Write)
}

type httpListener struct {
	server         *http.Server
	backlogChannel <-chan net.Conn
	closed         chan struct{}
}

func (h *httpListener) CloseChan() <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (h *httpListener) Close() error {
	//TODO implement me
	panic("implement me")
}

func (h *httpListener) Accept() (socket types.Socket, err error) {
	//TODO implement me
	panic("implement me")
}

func (h *httpListener) Addr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (h *httpListener) AcceptConn() (socket types.StreamConn, err error) {
	//TODO implement me
	panic("implement me")
}
