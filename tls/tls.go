//go:build !tls
// +build !tls

package tls

import (
	"crypto/tls"
	"github.com/hadi77ir/muxedsocket/utils"
	"net"
)

func ServerTLS(conn net.Conn, params any) net.Conn {
	var config *tls.Config
	if c, ok := params.(*tls.Config); ok {
		config = c
	}
	return tls.Server(conn, config)
}

func ClientTLS(conn net.Conn, params any) net.Conn {
	var config *tls.Config
	if c, ok := params.(*tls.Config); ok {
		config = c
	}
	return tls.Client(conn, config)
}

func ParseTLS(parameters utils.Parameters, isClient bool) (any, error) {
	config := &tls.Config{
		ServerName: GetSNIFromParams(parameters),
		NextProtos: GetNextProtosFromParams(parameters),
	}

	if isClient {
		verifierFunc, insecure, err := GetCertificatePinningAndInsecure(parameters)
		if err != nil {
			return nil, err
		}
		config.InsecureSkipVerify = insecure
		config.VerifyPeerCertificate = verifierFunc
	}

	if !isClient {
		clientCaPool, clientCaLen, err := LoadCertPoolFromParams(parameters, ParamClientCA)
		if err != nil {
			return nil, err
		}
		clientAuth := tls.NoClientCert
		if clientCaLen > 0 {
			clientAuth = tls.RequireAndVerifyClientCert
		}
		config.ClientCAs = clientCaPool
		config.ClientAuth = clientAuth
	}

	certs, err := LoadX509PairsFromParams(parameters)
	if err != nil {
		return nil, err
	}
	config.Certificates = certs

	return config, nil
}

func LoadX509PairsFromParams(parameters utils.Parameters) ([]tls.Certificate, error) {
	certs, keys, err := LoadX509PairsBytesFromParams(parameters)
	if err != nil {
		return nil, err
	}
	pairs := make([]tls.Certificate, len(keys))
	for i := 0; i < len(keys); i++ {
		pairs[i], err = tls.X509KeyPair(certs[i], keys[i])
		if err != nil {
			return nil, err
		}
	}
	return pairs, nil
}

func LoadX509PairFromParams(parameters utils.Parameters) (tls.Certificate, error) {
	cert, key, err := LoadX509PairBytesFromParams(parameters)
	if err != nil {
		var zero tls.Certificate
		return zero, err
	}
	return tls.X509KeyPair(cert, key)
}
