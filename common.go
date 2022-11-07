package muxedsocket

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var SchemeNotSupported = errors.New("scheme not supported")

const defaultPacketTransport = "udp"
const defaultStreamTransport = "tcp"
const defaultMuxer = "yamux"

// remove an element at given index. don't use if order matters in slice.
// from https://stackoverflow.com/a/37335777
func remove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func findAndRemove[T any](keys []string, registry *Registry[T]) ([]string, string, T, bool) {
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		if value, found := registry.Get(key); found {
			keys = remove(keys, i)
			return keys, key, value, true
		}
	}
	var none T
	return keys, "", none, false
}
func findAndRemoveTLS(parts []string) ([]string, bool) {
	for i := 0; i < len(parts); i++ {
		if parts[i] == "secure" || parts[i] == "tls" {
			return remove(parts, i), true
		}
	}
	return parts, false
}
func GetCommonParamsFromURL(addr *url.URL) CommonParams {
	keepalive, _ := time.ParseDuration(addr.Query().Get("keepalive"))
	timeout, _ := time.ParseDuration(addr.Query().Get("timeout"))
	return CommonParams{KeepalivePeriod: keepalive, MaxIdleTimeout: timeout}
}

func GetAddrByTransportType(transportKey string, addr string) (NetAddrPort, error) {
	sep := strings.LastIndex(addr, ":")
	var ipStr string
	var portStr string
	if sep != -1 {
		ipStr = addr[:sep]
		portStr = addr[sep+1:]
	} else {
		ipStr = addr
		portStr = "0"
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, net.InvalidAddrError("invalid ip address")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	if transportKey == "udp" {
		return &net.UDPAddr{IP: ip, Port: port}, nil
	}
	if transportKey == "tcp" {
		return &net.TCPAddr{IP: ip, Port: port}, nil
	}

	return &netAddr{
		host:    addr,
		port:    uint16(port),
		network: transportKey,
	}, nil
}
