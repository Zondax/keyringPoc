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
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/hashicorp/go-plugin"

	plugin2 "github.com/zondax/keyringPoc/keyring/grpc"
	keyring2 "github.com/zondax/keyringPoc/keyring/types"
)

const (
	backendId = "memoryGo"
	v         = "0.0.1"
)

type memKeyring struct {
	cdc codec.Codec
	db  db
}

func newMemKeyring() *memKeyring {
	registry := codectypes.NewInterfaceRegistry()
	cryptoCodec.RegisterInterfaces(registry)
	return &memKeyring{
		cdc: codec.NewProtoCodec(registry),
		db:  newDb(),
	}
}

func (k memKeyring) Backend(r *keyring2.BackendRequest) (*keyring2.BackendResponse, error) {
	return &keyring2.BackendResponse{Backend: fmt.Sprintf(`%s - %s`, backendId, v)}, nil
}

func (k memKeyring) Key(r *keyring2.KeyRequest) (*keyring2.KeyResponse, error) {
	item, err := k.db.get(fmt.Sprintf("%s.%s", r.Uid, "info"))
	if err != nil {
		return nil, err
	}
	return &keyring2.KeyResponse{
		Key: item,
	}, nil
}

func (k memKeyring) NewAccount(r *keyring2.NewAccountRequest) (*keyring2.NewAccountResponse, error) {
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

	k.db.save(fmt.Sprintf("%s.%s", r.Uid, "info"), serializedRecord)
	k.db.save(fmt.Sprintf("%s.%s", hex.EncodeToString(addr.Bytes()), "address"), []byte(fmt.Sprintf("%s.%s", r.Uid, "info")))
	return &keyring2.NewAccountResponse{Record: record}, nil
}

func extractPrivKeyFromLocal(rl *cosmosKeyring.Record_Local) (cryptotypes.PrivKey, error) {
	if rl.PrivKey == nil {
		return nil, errors.New("no priv key")
	}

	priv, ok := rl.PrivKey.GetCachedValue().(cryptotypes.PrivKey)
	if !ok {
		return nil, errors.New("no cached value")
	}

	return priv, nil
}

func (k memKeyring) Sign(r *keyring2.NewSignRequest) (*keyring2.NewSignResponse, error) {
	item, err := k.db.get(fmt.Sprintf("%s.%s", r.Uid, "info"))
	if err != nil {
		return nil, err
	}

	record := new(cosmosKeyring.Record)
	err = k.cdc.Unmarshal(item, record)
	if err != nil {
		return nil, err
	}
	switch {
	case record.GetLocal() != nil:
		priv, err := extractPrivKeyFromLocal(record.GetLocal())
		if err != nil {
			return nil, err
		}

		sig, err := priv.Sign(r.GetMsg())
		if err != nil {
			return nil, err
		}

		privKey, err := codectypes.NewAnyWithValue(priv.PubKey())
		if err != nil {
			return nil, errors.New("sdfd")
		}
		return &keyring2.NewSignResponse{
			Msg:    sig,
			PubKey: privKey,
		}, nil

	default:
		_, err := record.GetPubKey()
		if err != nil {
			return nil, err
		}
		return nil, errors.New("cannot sign with offline keys")
	}

}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin2.Handshake,
		Plugins: map[string]plugin.Plugin{
			"keyring": &plugin2.KeyringGRPC{Impl: newMemKeyring()},
		},
		// A non-nil value here enables gRPC serving for this keyStore...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
