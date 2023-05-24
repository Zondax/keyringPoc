package main

import (
	"context"
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/go-bip39"
	"github.com/zondax/keyringPoc/keyring"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	registry := codectypes.NewInterfaceRegistry()
	cryptoCodec.RegisterInterfaces(registry)
	//cdc := codec.NewProtoCodec(registry)

	fmt.Println("HELLO")
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := keyring.NewKeyringServiceClient(conn)

	r, err := client.Backend(context.Background(), &keyring.BackendRequest{})
	if err != nil {
		panic(err)
	}

	fmt.Println(r.Backend)
	mnemonic, err := newMnemonic()
	accountRequest := &keyring.NewAccountRequest{
		Uid:             "test",
		Mnemonic:        mnemonic,
		Bip39Passphrase: "",
		Hdpath:          hd.CreateHDPath(118, 0, 0).String(),
	}
	_, err = client.NewAccount(context.Background(), accountRequest)
	if err != nil {
		panic(err)
	}

	key, err := client.Key(context.Background(), &keyring.KeyRequest{Uid: "test"})
	//key.Record.PubKey.UnmarshalJSON()
	//rec := new(cosmosKeyring.Record)
	//err = rec.Unmarshal(key.Record.PubKey.Value)
	//err = cdc.Unmarshal(key.Record, rec)

	fmt.Println(key.Record.GetAddress())
}

func newMnemonic() (string, error) {
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}
