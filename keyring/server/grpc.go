package server

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cosmosKeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/zondax/keyringPoc/keyring"
)

type ApiServer struct {
	keyring.UnimplementedKeyringServiceServer
	cdc codec.Codec
	db  map[string][]byte
}

func NewApiServer() *ApiServer {
	registry := codectypes.NewInterfaceRegistry()
	cryptoCodec.RegisterInterfaces(registry)
	return &ApiServer{
		cdc: codec.NewProtoCodec(registry),
		db:  make(map[string][]byte),
	}
}

func (s ApiServer) Backend(ctx context.Context, r *keyring.BackendRequest) (*keyring.BackendResponse, error) {
	return &keyring.BackendResponse{Backend: "test"}, nil
}

func (s ApiServer) Key(ctx context.Context, r *keyring.KeyRequest) (*keyring.KeyResponse, error) {
	item, ok := s.db[fmt.Sprintf("%s.%s", r.Uid, "info")]
	if !ok {
		return nil, errors.New("key not found")
	}
	if len(item) == 0 {
		return nil, errors.New("error key")
	}

	k := new(cosmosKeyring.Record)
	err := s.cdc.Unmarshal(item, k)
	if err == nil {
		return &keyring.KeyResponse{
			Record: k,
		}, nil
	}

	return nil, errors.New("COULD NOT DECODE")
}

func (s ApiServer) NewAccount(ctx context.Context, r *keyring.NewAccountRequest) (*keyring.NewAccountResponse, error) {

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

	serializedRecord, err := s.cdc.Marshal(record)
	if err != nil {
		return nil, err
	}

	s.db[fmt.Sprintf("%s.%s", r.Uid, "info")] = serializedRecord
	s.db[fmt.Sprintf("%s.%s", hex.EncodeToString(addr.Bytes()), "address")] = []byte(fmt.Sprintf("%s.%s", r.Uid, "info"))
	return &keyring.NewAccountResponse{Record: record}, nil
}
