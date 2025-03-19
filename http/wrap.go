package http

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/hadi77ir/muxedsocket/basics/stream"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
	"golang.org/x/net/http2"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var ErrHostNotDefined = errors.New("host is required")

func WrapHttpClient(dialFunc types.StreamDialFunc, parameters utils.Parameters) (types.StreamDialFunc, error) {
	dialer, err := wrapStandardHttpClient(func() (net.Conn, error) {
		return dialFunc()
	}, parameters)
	if err != nil {
		return nil, err
	}
	return stream.WrapDialer(dialer), nil
}

func wrapStandardHttpClient(dialFunc stream.StandardPrimedDialFunc, parameters utils.Parameters) (stream.StandardPrimedDialFunc, error) {
	transport, scheme, err := CreateTransport(dialFunc, parameters)
	if err != nil {
		return nil, err
	}
	requestConstructor, err := CreateRequestConstructor(parameters, scheme)
	if err != nil {
		return nil, err
	}
	return func() (net.Conn, error) {
		bodyReader, bodyWriter := net.Pipe()
		request, err := requestConstructor(bodyReader)
		if err != nil {
			return nil, err
		}
		response, err := transport.RoundTrip(request)
		return stream.WrapFifoConn(response.Body, bodyWriter), nil
	}, nil
}

type RequestConstructorFunc func(reqBody io.ReadCloser) (*http.Request, error)

func CreateTransport(dialFunc stream.StandardPrimedDialFunc, parameters utils.Parameters) (roundTripper http.RoundTripper, scheme string, err error) {
	dialer := WrapDialerWithContext(func(network, addr string) (net.Conn, error) {
		return dialFunc()
	})

	forceH2, scheme := GetProtocolFromParameters(parameters)
	if forceH2 {
		roundTripper = &http2.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				return dialer(ctx, network, addr)
			},
			DisableCompression: !utils.BoolFromParameters(parameters, "compression", true),
			AllowHTTP:          true,
		}
		return
	}

	roundTripper = &http.Transport{
		DialTLSContext:      dialer,
		DialContext:         dialer,
		ForceAttemptHTTP2:   false,
		DisableCompression:  !utils.BoolFromParameters(parameters, "compression", true),
		DisableKeepAlives:   !utils.BoolFromParameters(parameters, "keepalive", true),
		MaxIdleConns:        utils.IntegerFromParameters(parameters, "idle", 0),
		MaxIdleConnsPerHost: utils.IntegerFromParameters(parameters, "idleperhost", http.DefaultMaxIdleConnsPerHost),
	}
	return
}

const (
	ProtoHTTP  = "http"
	ProtoHTTPS = "https"
)

func GetProtocolFromParameters(parameters utils.Parameters) (bool, string) {
	switch utils.StringFromParameters(parameters, "proto", "http") {
	case "https":
		// HTTP/1.1 over TLS. doesn't really differ from cleartext connection.
		return false, ProtoHTTPS
	case "http":
		// HTTP/1.1 Cleartext. doesn't really differ from TLS connection.
		return false, ProtoHTTP
	case "h2":
		// HTTP/2. There are differences between this and cleartext connection.
		return true, ProtoHTTPS
	case "h2c":
		// HTTP/2 Cleartext. There are differences between this and secure connection.
		return true, ProtoHTTP
	}
	// fallback to HTTP/1.1 cleartext.
	return false, ProtoHTTP
}

type dialerResult struct {
	conn net.Conn
	err  error
}

type DialTLSContextFunc func(ctx context.Context, network string, addr string, config *tls.Config) (net.Conn, error)
type DialContextFunc func(ctx context.Context, network string, addr string) (net.Conn, error)

func WrapDialerWithContext(dialer stream.StandardDialFunc) DialContextFunc {
	return func(ctx context.Context, network string, addr string) (net.Conn, error) {
		channel := make(chan dialerResult, 1)
		go func() {
			conn, err := dialer(network, addr)
			channel <- dialerResult{conn: conn, err: err}
			close(channel)
		}()
		select {
		case result, _ := <-channel:
			return result.conn, result.err
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			return nil, &net.OpError{Op: "dial", Err: context.DeadlineExceeded, Net: network, Source: types.EmptyAddr("src"), Addr: types.EmptyAddr("dst")}
		}
	}
}

func CreateRequestConstructor(parameters utils.Parameters, scheme string) (RequestConstructorFunc, error) {
	remoteUrl := createURLFromParameters(parameters, scheme)
	remoteUrlStr := remoteUrl.String()
	method := strings.ToUpper(utils.StringFromParameters(parameters, "method", "GET"))
	switch method {
	case "CONNECT":
		connectUrl, err := createConnectURLFromParameters(scheme, remoteUrl.Host)
		if err != nil {
			return nil, err
		}
		connectUrlStr := connectUrl.String()
		return func(reqBody io.ReadCloser) (*http.Request, error) {
			req := http.NewRequest(http.MethodConnect, connectUrlStr, reqBody)
			req.Header.Add("Proxy-Authorization")
		}, nil
	default:
		return func(reqBody io.ReadCloser) (*http.Request, error) {
			return http.NewRequest(method, remoteUrlStr, reqBody)
		}, nil
	}
}

func createConnectURLFromParameters(scheme, host string) (*url.URL, error) {
	connectUrl := &url.URL{
		Scheme: scheme,
		Host:   host,
	}
	if connectUrl.Host == "" {
		return nil, ErrHostNotDefined
	}
	return connectUrl, nil
}

func createURLFromParameters(parameters utils.Parameters, scheme string) *url.URL {
	userInfo := userInfoFromParameters(parameters)
	requestURL := &url.URL{
		Host: utils.StringFromParameters(parameters, "host", ""),
		// Scheme should be set to "https", regardless of whether our connection is encrypted or not.
		// This will prevent HTTP client to introduce TLS itself.
		Scheme:   scheme,
		User:     userInfo,
		Path:     utils.StringFromParameters(parameters, "path", ""),
		RawQuery: utils.StringFromParameters(parameters, "query", ""),
	}
	if !strings.HasPrefix(requestURL.Path, "/") {
		requestURL.Path = "/" + requestURL.Path
	}
	return requestURL
}

func userInfoFromParameters(parameters utils.Parameters) *url.Userinfo {
	if username, found := parameters.Get("username"); found {
		if password, found := parameters.Get("password"); found {
			return url.UserPassword(username, password)
		}
		return url.User(username)
	}
	return nil
}
