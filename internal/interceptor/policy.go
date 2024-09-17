package interceptor

import (
	"context"
	"errors"

	"github.com/8thgencore/microservice-chat/internal/client/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Client contains client connection with authentication service.
type Client struct {
	Client rpc.AuthClient
}

// PolicyInterceptor is used for authorization.
func (c *Client) PolicyInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	err := c.Client.Check(metadata.NewOutgoingContext(ctx, md), info.FullMethod)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
