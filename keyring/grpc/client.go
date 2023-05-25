package grpc

import (
	"context"
	"github.com/zondax/keyringPoc/keyring/types"
)

type Client struct{ Client types.KeyringServiceClient }

func (c *Client) Backend(r *types.BackendRequest) (*types.BackendResponse, error) {
	return c.Client.Backend(context.Background(), r)
}

func (c *Client) Key(r *types.KeyRequest) (*types.KeyResponse, error) {
	return c.Client.Key(context.Background(), r)
}

func (c *Client) NewAccount(r *types.NewAccountRequest) (*types.NewAccountResponse, error) {
	return c.Client.NewAccount(context.Background(), r)
}
