package main

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cosmosKeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
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
		mnemonic string
		hdPath   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "local priv key",
			args: args{
				mnemonic: "once capable omit cancel ghost mobile mean surface neither tissue life huge knock rebuild work enemy avoid bargain swarm paper comic follow blade tribe",
				hdPath:   hdPath,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			derivedPriv, err := hd.Secp256k1.Derive()(tt.args.mnemonic, "", tt.args.hdPath)
			require.NoError(t, err)
			require.NotEmpty(t, derivedPriv)
			priv := hd.Secp256k1.Generate()(derivedPriv)

			v, err := codectypes.NewAnyWithValue(priv)
			require.NoError(t, err)
			require.NotNil(t, v)

			recordLocal := &cosmosKeyring.Record_Local{v}
			recordLocalItem := &cosmosKeyring.Record_Local_{recordLocal}

			pk, err := codectypes.NewAnyWithValue(priv.PubKey())
			require.NoError(t, err)
			require.NotNil(t, pk)

			record := &cosmosKeyring.Record{"test", pk, recordLocalItem}
			got, err := extractPrivKeyFromLocal(record.GetLocal())
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, priv, got)
		})
	}
}

func Test_file_Backend(t *testing.T) {
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

func Test_file_Key(t *testing.T) {
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

func Test_file_NewAccount(t *testing.T) {
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

func Test_file_Sign(t *testing.T) {
	type args struct {
		r *keyring2.SignRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sign message",
			args: args{r: &keyring2.SignRequest{
				Uid:      "test",
				Msg:      []byte("this is a string"),
				SignMode: 0,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := newFileKeyring()
			_, err := k.NewAccount(&keyring2.NewAccountRequest{
				Uid:             tt.args.r.Uid,
				Mnemonic:        "once capable omit cancel ghost mobile mean surface neither tissue life huge knock rebuild work enemy avoid bargain swarm paper comic follow blade tribe",
				Bip39Passphrase: "",
				Hdpath:          "m/44'/118'/0'/0/0",
			})
			require.NoError(t, err)
			got, err := k.Sign(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			require.NotEmpty(t, got)
		})
	}
}

func Test_newfile(t *testing.T) {
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
			got := newFileKeyring()
			require.Equal(t, got, tt.want)
		})
	}
}

func Test_fileKeyring_SaveOffline(t *testing.T) {
	tests := []struct {
		name    string
		keyName string
		want    *keyring2.SaveOfflineResponse
		wantErr bool
	}{
		{
			name:    "save offline",
			keyName: "offlineTesting",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := newFileKeyring()
			pubKey := secp256k1.GenPrivKey().PubKey()
			pb, err := codectypes.NewAnyWithValue(pubKey)
			require.NoError(t, err)
			require.NotNil(t, pb)
			got, err := k.SaveOffline(&keyring2.SaveOfflineRequest{
				Uid:    tt.keyName,
				PubKey: pb,
			})
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
