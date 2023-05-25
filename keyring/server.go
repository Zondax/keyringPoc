package keyring

import (
	"context"
)

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	UnimplementedKeyringServiceServer
	// This is the real implementation
	Impl Keyring
}

func (s GRPCServer) Backend(ctx context.Context, r *BackendRequest) (*BackendResponse, error) {
	return s.Impl.Backend(r)
}

func (s GRPCServer) Key(ctx context.Context, r *KeyRequest) (*KeyResponse, error) {
	return s.Impl.Key(r)
}

func (s GRPCServer) NewAccount(ctx context.Context, r *NewAccountRequest) (*NewAccountResponse, error) {
	return s.Impl.NewAccount(r)
}
