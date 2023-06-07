package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/go-bip39"

	"github.com/zondax/keyringPoc/keyring/keyStore"
)

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

func main() {

	registry := codectypes.NewInterfaceRegistry()
	cryptoCodec.RegisterInterfaces(registry)

	k := keyStore.NewKeyring("./build/goFile", codec.NewProtoCodec(registry))
	defer k.Close()

	// Get the backend of the keyStore
	back := k.Backend()
	fmt.Println(back)

	// Create a new key
	//mnemonic, _ := newMnemonic()
	r, err := k.NewAccount(
		"test",
		"spare august spell toilet open wonder coffee tiger prepare size option talent citizen hungry vote swarm embark citizen hedgehog age giggle foster flat police",
		"",
		hd.CreateHDPath(118, 0, 0).String(),
		nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)

	record, err := k.Key("test")
	if err != nil {
		panic(err)
	}
	addr, err := record.GetAddress()
	fmt.Println(addr.String())
	fmt.Println(record.GetAddress())
	//msg := []byte("THIS IS A MESSAGE")
	//s, pk, err := k.Sign("test", msg)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(s))
	//fmt.Println(pk)
}
