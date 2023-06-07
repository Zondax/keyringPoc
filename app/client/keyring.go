package client

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcodec "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/zondax/keyringPoc/keyring/keyStore"
)

var availablePlugins = map[string]string{
	"goFile": "./build/goFile",
	"pyFile": "python3 plugin/pyFile/pyFile.py",
}

func GetCodec() (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	ir := codectypes.NewInterfaceRegistry()
	authcodec.RegisterInterfaces(ir)
	bankcodec.RegisterInterfaces(ir)
	cryptoCodec.RegisterInterfaces(ir)

	cdc := codec.NewProtoCodec(ir)
	return cdc, ir
}

func GetKeyring(plugin string) (*keyStore.PluginsKeyStore, error) {
	cdc, _ := GetCodec()
	p, ok := availablePlugins[plugin]
	if !ok {
		return nil, errors.New("not available plugin")
	}
	k := keyStore.NewKeyring(p, cdc)
	return k, nil
}
