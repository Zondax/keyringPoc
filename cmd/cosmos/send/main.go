package main

import (
	"fmt"
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
	hdPath = "m/44'/118'/0'/0/0"
)

func main() {
	c, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	cli, err := client.NewClientFromNode("https://cosmos-rpc.polkachu.com:443")
	if err != nil {
		panic(err)
	}

	ir := codectypes.NewInterfaceRegistry()
	authcodec.RegisterInterfaces(ir)
	bankcodec.RegisterInterfaces(ir)
	cryptoCodec.RegisterInterfaces(ir)

	cdc := codec.NewProtoCodec(ir)

	ctx := client.Context{}.
		WithCodec(cdc).
		WithInterfaceRegistry(ir).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithTxConfig(tx.NewTxConfig(cdc, tx.DefaultSignModes)).
		WithViper(""). // In simapp, we don't use any prefix for env variables.
		WithClient(cli).
		WithChainID("cosmoshub-4").
		WithFromName("test").
		WithBroadcastMode("sync")
	fmt.Println(ctx)

	//kr := keyring.NewInMemory(cdc)
	//key, err := kr.NewAccount("test", mnemonic, "", hdPath, hd.Secp256k1)
	//if err != nil {
	//	panic(err)
	//}
	//add, err := key.GetAddress()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(add.String())

	ks := keyStore.NewKeyring("./build/memoryGo", cdc)
	defer ks.Close()

	_, err = ks.NewAccount("test", c.Mnemonic, "", hdPath, hd.Secp256k1)
	if err != nil {
		panic(err)
	}
	key, err := ks.Key("test")
	add, err := key.GetAddress()
	if err != nil {
		panic(err)
	}
	fmt.Println(add.String())

	//amount := "10000uatom"
	coins, err := sdk.ParseCoinsNormalized("10000uatom")
	if err != nil {
		panic(err)
	}
	msg := &bankcodec.MsgSend{
		FromAddress: add.String(),
		ToAddress:   "cosmos1hw0lawgqtm0segnt34yuh63glujwv9kr6r0evp",
		Amount:      coins,
	}

	if err := msg.ValidateBasic(); err != nil {
		panic(err)
	}

	txf := txClient.Factory{}.WithGas(350000).
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
	//WithAccountNumber(1685823).
	//WithSequence(4).

	if err != nil {
		panic(err)
	}
	ctx = ctx.WithFromAddress(add)
	err = txClient.GenerateOrBroadcastTxWithFactory(ctx, txf, msg)
	if err != nil {
		panic(err)
	}
}
