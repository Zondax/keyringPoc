package keyring

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/hashicorp/go-plugin"
	"os"
	"os/exec"
)

type PluginsKeyring struct {
	client  *plugin.Client
	backEnd Keyring
	cdc     codec.Codec
}

func NewKeyring(path string, cdc codec.Codec) *PluginsKeyring {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  Handshake,
		Plugins:          PluginMap,
		Cmd:              exec.Command("sh", "-c", path),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	// Request the plugin
	raw, err := rpcClient.Dispense("keyring")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	return &PluginsKeyring{
		client:  client,
		backEnd: raw.(Keyring),
		cdc:     cdc,
	}
}

func (k PluginsKeyring) Close() {
	k.client.Kill()
}

func (k PluginsKeyring) Backend() string {
	backend, err := k.backEnd.Backend(&BackendRequest{})
	if err != nil {
		return ""
	}
	return backend.Backend
}

func (k PluginsKeyring) List() ([]*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) SupportedAlgorithms() (keyring.SigningAlgoList, keyring.SigningAlgoList) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) Key(uid string) (*keyring.Record, error) {
	r, err := k.backEnd.Key(&KeyRequest{Uid: uid})
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

func (k PluginsKeyring) KeyByAddress(address types.Address) (*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) Delete(uid string) error {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) DeleteByAddress(address types.Address) error {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) Rename(from string, to string) error {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) NewMnemonic(uid string, language keyring.Language, hdPath, bip39Passphrase string, algo keyring.SignatureAlgo) (*keyring.Record, string, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) NewAccount(uid, mnemonic, bip39Passphrase, hdPath string, algo keyring.SignatureAlgo) (*keyring.Record, error) {
	res, err := k.backEnd.NewAccount(&NewAccountRequest{
		Uid:      uid,
		Mnemonic: mnemonic,
		Hdpath:   hdPath,
	})
	if err != nil {
		return nil, err
	}
	return res.Record, nil
}

func (k PluginsKeyring) SaveLedgerKey(uid string, algo keyring.SignatureAlgo, hrp string, coinType, account, index uint32) (*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) SaveOfflineKey(uid string, pubkey types.PubKey) (*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) SaveMultisig(uid string, pubkey types.PubKey) (*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) Sign(uid string, msg []byte) ([]byte, types.PubKey, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) SignByAddress(address types.Address, msg []byte) ([]byte, types.PubKey, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) ImportPrivKey(uid, armor, passphrase string) error {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) ImportPubKey(uid string, armor string) error {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) ExportPubKeyArmor(uid string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) ExportPubKeyArmorByAddress(address types.Address) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) ExportPrivKeyArmor(uid, encryptPassphrase string) (armor string, err error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) ExportPrivKeyArmorByAddress(address types.Address, encryptPassphrase string) (armor string, err error) {
	//TODO implement me
	panic("implement me")
}

func (k PluginsKeyring) MigrateAll() ([]*keyring.Record, error) {
	//TODO implement me
	panic("implement me")
}
