package main

import (
	"errors"
	"fmt"
	"log"
	"os"

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
	backendId = "goFile"
)

type fileKeyring struct {
	cdc codec.Codec
	dir string
}

func newFileKeyring() *fileKeyring {
	registry := codectypes.NewInterfaceRegistry()
	cryptoCodec.RegisterInterfaces(registry)
	dir := os.TempDir() + "goPluginKeyring"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	return &fileKeyring{
		cdc: codec.NewProtoCodec(registry),
		dir: dir,
	}
}

func (k fileKeyring) Backend(r *keyring2.BackendRequest) (*keyring2.BackendResponse, error) {
	return &keyring2.BackendResponse{Backend: backendId}, nil
}

func (k fileKeyring) Key(r *keyring2.KeyRequest) (*keyring2.KeyResponse, error) {
	item, err := os.ReadFile(k.dir + fmt.Sprintf("/%s.%s", r.Uid, "info"))
	if err != nil {
		return nil, err
	}
	return &keyring2.KeyResponse{
		Key: item,
	}, nil
}

func (k fileKeyring) NewAccount(r *keyring2.NewAccountRequest) (*keyring2.NewAccountResponse, error) {
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

	serializedRecord, err := k.cdc.Marshal(record)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(k.dir+fmt.Sprintf("/%s.%s", r.Uid, "info"), serializedRecord, 0644)
	if err != nil {
		return nil, err
	}
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

func (k fileKeyring) Sign(r *keyring2.SignRequest) (*keyring2.SignResponse, error) {
	item, err := os.ReadFile(k.dir + fmt.Sprintf("/%s.%s", r.Uid, "info"))
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

		//pubKey, err := codectypes.NewAnyWithValue(priv.PubKey())
		if err != nil {
			return nil, errors.New("sdfd")
		}
		return &keyring2.SignResponse{
			Msg:    sig,
			Record: item,
		}, nil

	default:
		_, err := record.GetPubKey()
		if err != nil {
			return nil, err
		}
		return nil, errors.New("cannot sign with offline keys")
	}

}

func (k fileKeyring) SaveOffline(r *keyring2.SaveOfflineRequest) (*keyring2.SaveOfflineResponse, error) {
	recordOffline := &cosmosKeyring.Record_Offline{}
	recordOfflineItem := &cosmosKeyring.Record_Offline_{recordOffline}

	record := &cosmosKeyring.Record{r.Uid, r.PubKey, recordOfflineItem}
	// TODO: save file
	serializedRecord, err := k.cdc.Marshal(record)
	if err != nil {
		return nil, err
	}

	return &keyring2.SaveOfflineResponse{Record: serializedRecord}, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin2.Handshake,
		Plugins: map[string]plugin.Plugin{
			"keyring": &plugin2.KeyringGRPC{Impl: newFileKeyring()},
		},
		// A non-nil value here enables gRPC serving for this keyStore...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
