package keyring

import (
	"github.com/zondax/keyringPoc/keyring/types"
)

// PluginKeyring is the plugins interface
type PluginKeyring interface {
	Backend(*types.BackendRequest) (*types.BackendResponse, error)
	Key(*types.KeyRequest) (*types.KeyResponse, error)
	NewAccount(*types.NewAccountRequest) (*types.NewAccountResponse, error)
}
