package grpc

import (
	"context"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/zondax/keyringPoc/keyring"
	"github.com/zondax/keyringPoc/keyring/types"
)

// Handshake is a common handshake that is shared by keyStore and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN", // TODO
	MagicCookieValue: "hello",        // TODO
}

var PluginMap = map[string]plugin.Plugin{
	"keyring": &KeyringGRPC{},
}

type KeyringGRPC struct {
	plugin.Plugin
	Impl keyring.PluginKeyring
}

func (p *KeyringGRPC) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	types.RegisterKeyringServiceServer(s, &Server{Impl: p.Impl})
	return nil
}

func (p *KeyringGRPC) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &Client{Client: types.NewKeyringServiceClient(c)}, nil
}
