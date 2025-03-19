package muxedsocket

import (
	"github.com/hadi77ir/muxedsocket/chaining"
	"github.com/hadi77ir/muxedsocket/types"
	"github.com/hadi77ir/muxedsocket/utils"
	"net"
	"net/url"
)

type ListenFuncCreator[TFunc any] func(addr string, parameters utils.Parameters) (TFunc, error)
type MuxListenFuncCreator ListenFuncCreator[types.MuxListenFunc]

func ConstructListenFuncCreatorWithRegistry(creators *Creators, scheme string) (MuxListenFuncCreator, error) {
	if solution, found := creators.AddrSolutions().Get(scheme); found && solution != nil {
		return solution.Server, nil
	}
	schemeParts := GetSchemeParts(scheme)
	chainer, err := chaining.CreateMuxLayersChainer(creators, GetDefaults(creators), schemeParts)
	if err != nil {
		return nil, err
	}
	return chainer.ConstructListenFunc, nil
}

func ConstructListenFuncCreator(scheme string) (MuxListenFuncCreator, error) {
	return ConstructListenFuncCreatorWithRegistry(creators, scheme)
}

func CreateListenFuncWithRegistry(creators *Creators, addr *url.URL) (types.MuxListenFunc, error) {
	if addr == nil {
		return nil, net.InvalidAddrError("address was nil")
	}
	creatorFunc, err := ConstructListenFuncCreatorWithRegistry(creators, addr.Scheme)
	if err != nil {
		return nil, err
	}
	return creatorFunc(addr.Host, utils.ParametersFromURL(addr.Query()))
}

func CreateListenFunc(addr *url.URL) (types.MuxListenFunc, error) {
	return CreateListenFuncWithRegistry(creators, addr)
}

func ListenURIWithRegistry(creators *Creators, uri string) (types.MuxedListener, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	listener, err := CreateListenFuncWithRegistry(creators, parsed)
	if err != nil {
		return nil, err
	}
	return listener()
}

func ListenURI(uri string) (types.MuxedListener, error) {
	return ListenURIWithRegistry(creators, uri)
}
