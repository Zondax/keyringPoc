package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cosmosKeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/go-bip39"
	"github.com/stretchr/testify/require"
	keyring2 "github.com/zondax/keyringPoc/keyring/types"
	"os"
	"reflect"
	"testing"
)

const (
	hdPath = "m/44'/118'/0'/0/0"
)

func getCodec() *codec.ProtoCodec {
	registry := codectypes.NewInterfaceRegistry()
	cryptoCodec.RegisterInterfaces(registry)
	return codec.NewProtoCodec(registry)
}

func newMnemonic() string {
	entropySeed, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropySeed)
	return mnemonic
}

func Test_extractPrivKeyFromLocal(t *testing.T) {
	type args struct {
		rl *cosmosKeyring.Record_Local
	}
	tests := []struct {
		name    string
		args    args
		want    types.PrivKey
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractPrivKeyFromLocal(tt.args.rl)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractPrivKeyFromLocal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractPrivKeyFromLocal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memKeyring_Backend(t *testing.T) {
	type fields struct {
		cdc codec.Codec
	}
	type args struct {
		r *keyring2.BackendRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keyring2.BackendResponse
		wantErr bool
	}{
		{
			name: "backend response",
			fields: fields{
				cdc: getCodec(),
			},
			args:    args{r: &keyring2.BackendRequest{}},
			want:    &keyring2.BackendResponse{Backend: backendId},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := fileKeyring{
				cdc: tt.fields.cdc,
			}
			got, err := k.Backend(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Backend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Backend() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memKeyring_Key(t *testing.T) {
	type fields struct {
		cdc codec.Codec
		db  map[string][]byte
	}
	type args struct {
		r *keyring2.KeyRequest
	}
	tests := []struct {
		name      string
		args      args
		createKey string
		got       *cosmosKeyring.Record
		wantErr   bool
	}{
		{
			name:      "get key",
			args:      args{r: &keyring2.KeyRequest{Uid: "test"}},
			createKey: "test",
			got:       nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := newFileKeyring()
			if tt.createKey != "" {
				_, err := k.NewAccount(&keyring2.NewAccountRequest{
					Uid:             tt.args.r.Uid,
					Mnemonic:        newMnemonic(),
					Bip39Passphrase: "",
					Hdpath:          hdPath,
				})
				require.NoError(t, err)
			}
			_, err := k.Key(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Key() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_memKeyring_NewAccount(t *testing.T) {
	type fields struct {
		cdc codec.Codec
	}
	type args struct {
		r *keyring2.NewAccountRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "create key",
			fields: fields{
				cdc: getCodec(),
			},
			args: args{r: &keyring2.NewAccountRequest{
				Uid:             "testNewAccount",
				Mnemonic:        "once capable omit cancel ghost mobile mean surface neither tissue life huge knock rebuild work enemy avoid bargain swarm paper comic follow blade tribe",
				Bip39Passphrase: "",
				Hdpath:          hdPath,
			}},
			want:    "cosmos15ddgpspw20ugppyl03mg7j3r0kcv05mj05xjfk",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := newFileKeyring()
			got, err := k.NewAccount(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			add, err := got.Record.GetAddress()
			require.NoError(t, err)
			if !reflect.DeepEqual(add.String(), tt.want) {
				t.Errorf("NewAccount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memKeyring_Sign(t *testing.T) {
	type fields struct {
		cdc codec.Codec
	}
	type args struct {
		r *keyring2.NewSignRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *keyring2.NewSignResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := newFileKeyring()
			got, err := k.Sign(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sign() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newMemKeyring(t *testing.T) {
	tests := []struct {
		name string
		want *fileKeyring
	}{
		{
			name: "create memory keyring",
			want: &fileKeyring{
				cdc: getCodec(),
				dir: os.TempDir() + "goPluginKeyring",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newFileKeyring(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newMemKeyring() = %v, want %v", got, tt.want)
			}
		})
	}
}
