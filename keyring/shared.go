package keyring

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

var PluginMap = map[string]plugin.Plugin{
	"keyring": &KeyringGRPC{},
}

// Keyring is the plugins interface
type Keyring interface {
	Backend(*BackendRequest) (*BackendResponse, error)
	Key(*KeyRequest) (*KeyResponse, error)
	NewAccount(*NewAccountRequest) (*NewAccountResponse, error)
}

type KeyringGRPC struct {
	plugin.Plugin
	Impl Keyring
}

func (p *KeyringGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterKeyringServiceServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *KeyringGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{Client: NewKeyringServiceClient(c)}, nil
}
