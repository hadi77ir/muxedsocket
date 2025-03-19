package muxedsocket

import (
	"github.com/hadi77ir/muxedsocket/chaining"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
	"net"
	"net/url"
)

type DialFuncCreator[TFunc any] func(addr string, parameters utils.Parameters) (TFunc, error)
type MuxDialFuncCreator DialFuncCreator[types.MuxDialFunc]

func ConstructDialFuncCreatorWithRegistry(creators *Creators, scheme string) (MuxDialFuncCreator, error) {
	if solution, found := creators.AddrSolutions().Get(scheme); found && solution != nil {
		return solution.Client, nil
	}
	schemeParts := GetSchemeParts(scheme)
	chainer, err := chaining.CreateMuxLayersChainer(creators, GetDefaults(creators), schemeParts)
	if err != nil {
		return nil, err
	}
	return chainer.ConstructDialFunc, nil
}
func ConstructDialFuncCreator(scheme string) (MuxDialFuncCreator, error) {
	return ConstructDialFuncCreatorWithRegistry(creators, scheme)
}

func CreateDialerWithRegistry(creators *Creators, addr *url.URL) (types.MuxDialFunc, error) {
	if addr == nil {
		return nil, net.InvalidAddrError("address was nil")
	}
	creatorFunc, err := ConstructDialFuncCreatorWithRegistry(creators, addr.Scheme)
	if err != nil {
		return nil, err
	}
	return creatorFunc(addr.Host, utils.ParametersFromURL(addr.Query()))
}

func CreateDialer(addr *url.URL) (types.MuxDialFunc, error) {
	return CreateDialerWithRegistry(creators, addr)
}

func DialURIWithRegistry(creators *Creators, uri string) (types.MuxedSocket, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	dialer, err := CreateDialerWithRegistry(creators, parsed)
	if err != nil {
		return nil, err
	}
	return dialer()
}

func DialURI(uri string) (types.MuxedSocket, error) {
	return DialURIWithRegistry(creators, uri)
}
