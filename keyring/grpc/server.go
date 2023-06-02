package grpc

import (
	"context"

	"github.com/zondax/keyringPoc/keyring"
	"github.com/zondax/keyringPoc/keyring/types"
)

// Here is the gRPC server that GRPCClient talks to.
type Server struct {
	types.UnimplementedKeyringServiceServer
	// This is the real implementation
	Impl keyring.PluginKeyring
}

func (s Server) Backend(ctx context.Context, r *types.BackendRequest) (*types.BackendResponse, error) {
	return s.Impl.Backend(r)
}

func (s Server) Key(ctx context.Context, r *types.KeyRequest) (*types.KeyResponse, error) {
	return s.Impl.Key(r)
}

func (s Server) NewAccount(ctx context.Context, r *types.NewAccountRequest) (*types.NewAccountResponse, error) {
	return s.Impl.NewAccount(r)
}

func (s Server) Sign(ctx context.Context, r *types.NewSignRequest) (*types.NewSignResponse, error) {
	return s.Impl.Sign(r)
}
