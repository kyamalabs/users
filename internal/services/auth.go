package services

import (
	"context"
	"fmt"
	"time"

	"github.com/kyamalabs/auth/pkg/middleware"

	"github.com/kyamalabs/users/internal/constants"

	"github.com/kyamalabs/auth/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type AuthGrpcService interface {
	VerifyAccessToken(*pb.VerifyAccessTokenRequest, string) (*pb.VerifyAccessTokenResponse, error)
}

type AuthServiceGrpcClient struct {
	conn                   *grpc.ClientConn
	client                 pb.AuthClient
	serviceAuthPrivateKeys []string
}

type requestMetadata struct {
	accessToken string
}

func NewAuthServiceGrpcClient(serverAddress string, serviceAuthPrivateKeys []string) (AuthGrpcService, error) {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("could not create client connection to auth gRPC service: %w", err)
	}

	client := pb.NewAuthClient(conn)

	return &AuthServiceGrpcClient{
		conn:                   conn,
		client:                 client,
		serviceAuthPrivateKeys: serviceAuthPrivateKeys,
	}, nil
}

func (c *AuthServiceGrpcClient) Close() error {
	return c.conn.Close()
}

var getServiceAuthenticationPayload = func(serviceAuthPrivateKeys []string) (string, error) {
	return middleware.GenerateServiceAuthenticationPayload(constants.ServiceName, serviceAuthPrivateKeys)
}

func (c *AuthServiceGrpcClient) withMetadata(ctx context.Context, reqMetadata *requestMetadata) (context.Context, error) {
	serviceAuthenticationPayload, err := getServiceAuthenticationPayload(c.serviceAuthPrivateKeys)
	if err != nil {
		return nil, fmt.Errorf("could not generate service authentication payload: %w", err)
	}

	md := metadata.Pairs(
		constants.XServiceAuthenticationHeader, serviceAuthenticationPayload,
	)

	if reqMetadata.accessToken != "" {
		md.Append(constants.AuthorizationHeader, fmt.Sprintf("%s %s", constants.AuthorizationBearer, reqMetadata.accessToken))
	}

	return metadata.NewOutgoingContext(ctx, md), nil
}

func (c *AuthServiceGrpcClient) VerifyAccessToken(payload *pb.VerifyAccessTokenRequest, accessToken string) (*pb.VerifyAccessTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx, err := c.withMetadata(ctx, &requestMetadata{accessToken: accessToken})
	if err != nil {
		return nil, fmt.Errorf("could not add metadata to verify access token request: %w", err)
	}

	return c.client.VerifyAccessToken(ctx, payload)
}
