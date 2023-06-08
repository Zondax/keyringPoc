package keyStore

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cosmosKeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcodec "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	goPlugin     = "../../build/goFile"
	pythonPlugin = "python3 ../../plugin/pyFile/pyFile.py"
)

func getCodec() *codec.ProtoCodec {
	ir := codectypes.NewInterfaceRegistry()
	authcodec.RegisterInterfaces(ir)
	bankcodec.RegisterInterfaces(ir)
	cryptoCodec.RegisterInterfaces(ir)

	cdc := codec.NewProtoCodec(ir)
	return cdc
}

var cdc = getCodec()

func TestPluginsKeyStore_Backend(t *testing.T) {
	tests := []struct {
		name string
		k    cosmosKeyring.Keyring
		want string
	}{
		{
			name: "goFile plugin",
			k:    NewKeyring(goPlugin, cdc),
			want: "goFile",
		},
		{
			name: "pyFile plugin",
			k:    NewKeyring(pythonPlugin, cdc),
			want: "pyFile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.Backend(); got != tt.want {
				t.Errorf("Backend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluginsKeyStore_NewAccount(t *testing.T) {
	type args struct {
		uid      string
		mnemonic string
		hdPath   string
		algo     cosmosKeyring.SignatureAlgo
	}
	tests := []struct {
		name    string
		k       cosmosKeyring.Keyring
		args    args
		wantErr bool
	}{
		{
			name: "go plugin",
			k:    NewKeyring("../../build/goFile", cdc),
			args: args{
				uid:      "test",
				mnemonic: "spare august spell toilet open wonder coffee tiger prepare size option talent citizen hungry vote swarm embark citizen hedgehog age giggle foster flat police",
				hdPath:   "m/44'/118'/0'/0/0",
				algo:     hd.Secp256k1,
			},
			wantErr: false,
		},
		{
			name: "python plugin",
			k:    NewKeyring(pythonPlugin, cdc),
			args: args{
				uid:      "test",
				mnemonic: "spare august spell toilet open wonder coffee tiger prepare size option talent citizen hungry vote swarm embark citizen hedgehog age giggle foster flat police",
				hdPath:   "m/44'/118'/0'/0/0",
				algo:     hd.Secp256k1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.NewAccount(tt.args.uid, tt.args.mnemonic, "", tt.args.hdPath, tt.args.algo)
			require.NoError(t, err)
			require.NotNil(t, got.PubKey.Value)
			require.NotNil(t, got.GetLocal().PrivKey.Value)
			ck := cosmosKeyring.NewInMemory(cdc)
			r, err := ck.NewAccount(tt.args.uid, tt.args.mnemonic, "", tt.args.hdPath, tt.args.algo)
			require.NoError(t, err)
			require.Equal(t, r.PubKey.Value, got.PubKey.Value)
			require.Equal(t, r.GetLocal().PrivKey.Value, got.GetLocal().PrivKey.Value)
		})
	}
}

func TestPluginsKeyStore_NewAccountWithBip(t *testing.T) {
	type args struct {
		uid             string
		mnemonic        string
		bip39Passphrase string
		hdPath          string
		algo            cosmosKeyring.SignatureAlgo
	}
	tests := []struct {
		name    string
		k       cosmosKeyring.Keyring
		args    args
		wantErr bool
	}{
		{
			name: "go plugin with bip39Passphrase",
			k:    NewKeyring(goPlugin, cdc),
			args: args{
				uid:             "testBip",
				mnemonic:        "spare august spell toilet open wonder coffee tiger prepare size option talent citizen hungry vote swarm embark citizen hedgehog age giggle foster flat police",
				bip39Passphrase: "check",
				hdPath:          "m/44'/118'/0'/0/0",
				algo:            hd.Secp256k1,
			},
			wantErr: false,
		},
		{
			name: "python plugin with bip39Passphrase",
			k:    NewKeyring(pythonPlugin, cdc),
			args: args{
				uid:             "testBip",
				mnemonic:        "spare august spell toilet open wonder coffee tiger prepare size option talent citizen hungry vote swarm embark citizen hedgehog age giggle foster flat police",
				bip39Passphrase: "check",
				hdPath:          "m/44'/118'/0'/0/0",
				algo:            hd.Secp256k1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.NewAccount(tt.args.uid, tt.args.mnemonic, tt.args.bip39Passphrase, tt.args.hdPath, tt.args.algo)
			require.NoError(t, err)
			require.NotNil(t, got.PubKey.Value)
			require.NotNil(t, got.GetLocal().PrivKey.Value)
			ck := cosmosKeyring.NewInMemory(cdc)
			r, err := ck.NewAccount(tt.args.uid, tt.args.mnemonic, tt.args.bip39Passphrase, tt.args.hdPath, tt.args.algo)
			require.NoError(t, err)
			require.Equal(t, r.PubKey.Value, got.PubKey.Value)
			require.Equal(t, r.GetLocal().PrivKey.Value, got.GetLocal().PrivKey.Value)
		})
	}
}

func TestPluginsKeyStore_Key(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name    string
		k       cosmosKeyring.Keyring
		args    args
		address string
		wantErr bool
	}{
		{
			name:    "get go testBip",
			k:       NewKeyring(goPlugin, cdc),
			args:    args{uid: "test"},
			address: "cosmos10rn5l5lezc5pst4jgxjkev974zjg7tfhhg96gf",
			wantErr: false,
		},
		{
			name:    "get python testBip",
			k:       NewKeyring(pythonPlugin, cdc),
			args:    args{uid: "test"},
			address: "cosmos10rn5l5lezc5pst4jgxjkev974zjg7tfhhg96gf",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.Key(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Key() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			add, err := got.GetAddress()
			require.NoError(t, err)
			require.Equal(t, tt.address, add.String())
		})
	}
}

func TestPluginsKeyStore_KeyWithBip(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name    string
		k       cosmosKeyring.Keyring
		args    args
		address string
		wantErr bool
	}{
		{
			name:    "get go testBip",
			k:       NewKeyring(goPlugin, cdc),
			args:    args{uid: "testBip"},
			address: "cosmos1qeneag8ath9ppmz3rhkh8fgpyf72ygv9rdh97u",
			wantErr: false,
		},
		{
			name:    "get python testBip",
			k:       NewKeyring(pythonPlugin, cdc),
			args:    args{uid: "testBip"},
			address: "cosmos1qeneag8ath9ppmz3rhkh8fgpyf72ygv9rdh97u",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.Key(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Key() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			add, err := got.GetAddress()
			require.NoError(t, err)
			require.Equal(t, tt.address, add.String())
		})
	}
}

func TestPluginsKeyStore_Sign(t *testing.T) {
	type args struct {
		uid string
		msg []byte
	}
	type cosmosKeyringArgs struct {
		mnemonic        string
		bip39Passphrase string
	}
	tests := []struct {
		name       string
		k          cosmosKeyring.Keyring
		args       args
		cosmosArgs cosmosKeyringArgs
	}{
		{
			name: "sign with go plugin",
			k:    NewKeyring(goPlugin, cdc),
			args: args{
				uid: "test",
				msg: []byte("MESSAGE TO SIGN"),
			},
			cosmosArgs: cosmosKeyringArgs{
				mnemonic:        "spare august spell toilet open wonder coffee tiger prepare size option talent citizen hungry vote swarm embark citizen hedgehog age giggle foster flat police",
				bip39Passphrase: "",
			},
		},
		{
			name: "sign with python plugin",
			k:    NewKeyring(goPlugin, cdc),
			args: args{
				uid: "test",
				msg: []byte("MESSAGE TO SIGN"),
			},
			cosmosArgs: cosmosKeyringArgs{
				mnemonic:        "spare august spell toilet open wonder coffee tiger prepare size option talent citizen hungry vote swarm embark citizen hedgehog age giggle foster flat police",
				bip39Passphrase: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := tt.k.Sign(tt.args.uid, tt.args.msg)
			require.NoError(t, err)
			require.NotNil(t, got)
			ck := cosmosKeyring.NewInMemory(cdc)
			_, err = ck.NewAccount(tt.args.uid,
				tt.cosmosArgs.mnemonic,
				tt.cosmosArgs.bip39Passphrase,
				"m/44'/118'/0'/0/0",
				hd.Secp256k1)
			require.NoError(t, err)
			cosmosSign, _, err := ck.Sign(tt.args.uid, tt.args.msg)
			require.NoError(t, err)
			require.Equal(t, cosmosSign, got)
		})
	}
}
