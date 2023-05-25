package keyStore

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cyptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	pocKeyring "github.com/zondax/keyringPoc/keyring"
	"github.com/zondax/keyringPoc/keyring/grpc"
	"github.com/zondax/keyringPoc/keyring/types"
	"os"
	"os/exec"
)

type PluginsKeyStore struct {
	client  *plugin.Client
	backEnd pocKeyring.PluginKeyring
	cdc     codec.Codec
}

func NewKeyring(path string, cdc codec.Codec) *PluginsKeyStore {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  grpc.Handshake,
		Plugins:          grpc.PluginMap,
		Cmd:              exec.Command("sh", "-c", path),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           hclog.New(&hclog.LoggerOptions{Level: hclog.Off}),
	})
	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	// Request the keyStore
	raw, err := rpcClient.Dispense("keyring")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	return &PluginsKeyStore{
		client:  client,
		backEnd: raw.(pocKeyring.PluginKeyring),
		cdc:     cdc,
	}
}

func (k PluginsKeyStore) Close() {
	k.client.Kill()
}

func (k PluginsKeyStore) Backend() string {
	backend, err := k.backEnd.Backend(&types.BackendRequest{})
	if err != nil {
		return ""
	}
	return backend.Backend
}

func (k PluginsKeyStore) List() ([]*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) SupportedAlgorithms() (keyring.SigningAlgoList, keyring.SigningAlgoList) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) Key(uid string) (*keyring.Record, error) {
	r, err := k.backEnd.Key(&types.KeyRequest{Uid: uid})
	if err != nil {
		return nil, err
	}
	record := new(keyring.Record)
	err = k.cdc.Unmarshal(r.Key, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (k PluginsKeyStore) KeyByAddress(address cyptoTypes.Address) (*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) Delete(uid string) error {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) DeleteByAddress(address cyptoTypes.Address) error {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) Rename(from string, to string) error {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) NewMnemonic(uid string, language keyring.Language, hdPath, bip39Passphrase string, algo keyring.SignatureAlgo) (*keyring.Record, string, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) NewAccount(uid, mnemonic, bip39Passphrase, hdPath string, algo keyring.SignatureAlgo) (*keyring.Record, error) {
	res, err := k.backEnd.NewAccount(&types.NewAccountRequest{
		Uid:      uid,
		Mnemonic: mnemonic,
		Hdpath:   hdPath,
	})
	if err != nil {
		return nil, err
	}
	return res.Record, nil
}

func (k PluginsKeyStore) SaveLedgerKey(uid string, algo keyring.SignatureAlgo, hrp string, coinType, account, index uint32) (*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) SaveOfflineKey(uid string, pubkey cyptoTypes.PubKey) (*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) SaveMultisig(uid string, pubkey cyptoTypes.PubKey) (*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) Sign(uid string, msg []byte) ([]byte, cyptoTypes.PubKey, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) SignByAddress(address cyptoTypes.Address, msg []byte) ([]byte, cyptoTypes.PubKey, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) ImportPrivKey(uid, armor, passphrase string) error {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) ImportPubKey(uid string, armor string) error {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) ExportPubKeyArmor(uid string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) ExportPubKeyArmorByAddress(address cyptoTypes.Address) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) ExportPrivKeyArmor(uid, encryptPassphrase string) (armor string, err error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) ExportPrivKeyArmorByAddress(address cyptoTypes.Address, encryptPassphrase string) (armor string, err error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyStore) MigrateAll() ([]*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}
