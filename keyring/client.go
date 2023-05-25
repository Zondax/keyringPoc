package keyring

import (
	"context"
)

type GRPCClient struct{ Client KeyringServiceClient }

func (m *GRPCClient) Backend(r *BackendRequest) (*BackendResponse, error) {
	return m.Client.Backend(context.Background(), r)
}

func (m *GRPCClient) Key(r *KeyRequest) (*KeyResponse, error) {
	return m.Client.Key(context.Background(), r)
}

func (m *GRPCClient) NewAccount(r *NewAccountRequest) (*NewAccountResponse, error) {
	return m.Client.NewAccount(context.Background(), r)
}
