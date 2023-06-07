package tx

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	client2 "github.com/zondax/keyringPoc/app/client"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcodec "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func getCodec() (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	ir := codectypes.NewInterfaceRegistry()
	authcodec.RegisterInterfaces(ir)
	bankcodec.RegisterInterfaces(ir)
	cryptoCodec.RegisterInterfaces(ir)

	cdc := codec.NewProtoCodec(ir)
	return cdc, ir
}

func ctx(keyName, endpoint string, cdc *codec.ProtoCodec, ir codectypes.InterfaceRegistry) client.Context {
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

func Send(uid, plugin, toAddress, amount, node string) {
	k, err := client2.GetKeyring(plugin)
	if err != nil {
		panic(err)
	}
	defer k.Close()

	key, err := k.Key(uid)
	if err != nil {
		panic(err)
	}
	add, err := key.GetAddress()
	if err != nil {
		panic(err)
	}
	coins, err := sdk.ParseCoinsNormalized(amount)
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

	cdc, ir := getCodec()
	ctx := ctx(uid, node, cdc, ir)
	txf := txFactory(ctx, k)
	if err != nil {
		panic(err)
	}
	ctx = ctx.WithFromAddress(add)
	err = txClient.GenerateOrBroadcastTxWithFactory(ctx, txf, msg)
	if err != nil {
		panic(err)
	}
}
