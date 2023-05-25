package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cosmosKeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/hashicorp/go-plugin"
	keyring2 "github.com/zondax/keyringPoc/keyring"
)

type keyring struct {
	cdc codec.Codec
	db  map[string][]byte
}

func NewKeyring() *keyring {
	registry := codectypes.NewInterfaceRegistry()
	cryptoCodec.RegisterInterfaces(registry)
	return &keyring{
		cdc: codec.NewProtoCodec(registry),
		db:  make(map[string][]byte),
	}
}

func (k keyring) Backend(r *keyring2.BackendRequest) (*keyring2.BackendResponse, error) {
	return &keyring2.BackendResponse{Backend: "memoryGo"}, nil
}

func (k keyring) Key(r *keyring2.KeyRequest) (*keyring2.KeyResponse, error) {
	item, ok := k.db[fmt.Sprintf("%s.%s", r.Uid, "info")]
	if !ok {
		return nil, errors.New("key not found")
	}
	if len(item) == 0 {
		return nil, errors.New("error key")
	}

	return &keyring2.KeyResponse{
		Key: item,
	}, nil
}

func (k keyring) NewAccount(r *keyring2.NewAccountRequest) (*keyring2.NewAccountResponse, error) {
	derivedPriv, err := hd.Secp256k1.Derive()(r.Mnemonic, r.Bip39Passphrase, r.Hdpath)
	if err != nil {
		return nil, err
	}
	priv := hd.Secp256k1.Generate()(derivedPriv)

	v, err := codectypes.NewAnyWithValue(priv)
	if err != nil {
		return nil, err
	}

	recordLocal := &cosmosKeyring.Record_Local{v}
	recordLocalItem := &cosmosKeyring.Record_Local_{recordLocal}

	pk, err := codectypes.NewAnyWithValue(priv.PubKey())
	if err != nil {
		return nil, err
	}

	record := &cosmosKeyring.Record{r.Uid, pk, recordLocalItem}
	addr, err := record.GetAddress()
	if err != nil {
		return nil, err
	}

	serializedRecord, err := k.cdc.Marshal(record)
	if err != nil {
		return nil, err
	}

	k.db[fmt.Sprintf("%s.%s", r.Uid, "info")] = serializedRecord
	k.db[fmt.Sprintf("%s.%s", hex.EncodeToString(addr.Bytes()), "address")] = []byte(fmt.Sprintf("%s.%s", r.Uid, "info"))
	return &keyring2.NewAccountResponse{Record: record}, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: keyring2.Handshake,
		Plugins: map[string]plugin.Plugin{
			"keyring": &keyring2.KeyringGRPC{Impl: NewKeyring()},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
