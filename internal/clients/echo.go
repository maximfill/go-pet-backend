package clients

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/maximfill/go-pet-backend/internal/clients/pb"
)

type EchoClient struct {
	client pb.EchoServiceClient
}

func NewEchoClient(conn grpc.ClientConnInterface) *EchoClient {
	return &EchoClient{
		client: pb.NewEchoServiceClient(conn),
	}
}

func (c *EchoClient) Echo(ctx context.Context, msg string) error { // реализация вызова
	_, err := c.client.UnaryEcho(ctx, &pb.EchoRequest{
		Message: msg,
	})
	return err
}
