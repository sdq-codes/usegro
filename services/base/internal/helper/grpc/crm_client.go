package grpcclient

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/usegro/proto/crm"
)

type CRMClient struct {
	client pb.CRMServiceClient
}

func NewCRMClient(address string) (*CRMClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &CRMClient{client: pb.NewCRMServiceClient(conn)}, nil
}

func (c *CRMClient) ListOrganizationsByUser(ctx context.Context, userID string) ([]*pb.CrmOrganization, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.ListOrganizationsByUser(ctx, &pb.ListOrganizationsByUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return resp.Organizations, nil
}
