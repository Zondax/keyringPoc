package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcodec "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/zondax/keyringPoc/keyring/keyStore"
)

const (
	hdPath       = "m/44'/118'/0'/0/0"
	toAddress    = "cosmos1hw0lawgqtm0segnt34yuh63glujwv9kr6r0evp"
	rpcEndpoint  = "https://cosmos-rpc.polkachu.com:443"
	goPlugin     = "./build/goFile"
	pythonPlugin = "python3 plugin/pyFile/pyFile.py"
	keyName      = "test"
)

func checkCosmosKeyring(mnemonic string, cdc *codec.ProtoCodec) {
	kr := keyring.NewInMemory(cdc)
	key, err := kr.NewAccount(keyName, mnemonic, "", hdPath, hd.Secp256k1)
	if err != nil {
		panic(err)
	}
	add, err := key.GetAddress()
	if err != nil {
		panic(err)
	}
	fmt.Println(add.String())
}

func getCodec() (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	ir := codectypes.NewInterfaceRegistry()
	authcodec.RegisterInterfaces(ir)
	bankcodec.RegisterInterfaces(ir)
	cryptoCodec.RegisterInterfaces(ir)

	cdc := codec.NewProtoCodec(ir)
	return cdc, ir
}

func ctx(endpoint string, cdc *codec.ProtoCodec, ir codectypes.InterfaceRegistry) client.Context {
	cli, err := client.NewClientFromNode(endpoint)
	if err != nil {
		panic(err)
	}
	return client.Context{}.
		WithCodec(cdc).
		WithInterfaceRegistry(ir).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithTxConfig(tx.NewTxConfig(cdc, tx.DefaultSignModes)).
		WithViper(""). // TODO understand
		WithClient(cli).
		WithChainID("cosmoshub-4").
		WithFromName(keyName).
		WithBroadcastMode("sync")
}

func txFactory(ctx client.Context, ks keyring.Keyring) txClient.Factory {
	return txClient.Factory{}.WithGas(350000).
		WithSimulateAndExecute(false).
		WithGasAdjustment(1).
		WithMemo("").
		WithTimeoutHeight(0).
		WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
		WithTxConfig(ctx.TxConfig).
		WithAccountRetriever(ctx.AccountRetriever).
		WithChainID("cosmoshub-4").
		WithFees("2500uatom").
		WithKeybase(ks)
}

func main() {
	c, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	cdc, ir := getCodec()
	ctx := ctx(rpcEndpoint, cdc, ir)

	ks := keyStore.NewKeyring(pythonPlugin, cdc)
	defer ks.Close()
	fmt.Printf("Using keyring with plugin: %s\n\n", ks.Backend())

	_, err = ks.NewAccount(keyName, c.Mnemonic, "", hdPath, hd.Secp256k1)
	if err != nil {
		panic(err)
	}
	key, err := ks.Key(keyName)
	add, err := key.GetAddress()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Address:\n\t%s\n", add.String())

	coins, err := sdk.ParseCoinsNormalized("10000uatom")
	if err != nil {
		panic(err)
	}
	msg := &bankcodec.MsgSend{
		FromAddress: add.String(),
		ToAddress:   toAddress,
		Amount:      coins,
	}
	if err := msg.ValidateBasic(); err != nil {
		panic(err)
	}
	fmt.Printf("About to send:\n\t%v\n\n", msg)

	txf := txFactory(ctx, ks)
	if err != nil {
		panic(err)
	}
	ctx = ctx.WithFromAddress(add)
	err = txClient.GenerateOrBroadcastTxWithFactory(ctx, txf, msg)
	if err != nil {
		panic(err)
	}
}
