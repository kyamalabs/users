package services

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/kyamalabs/proto/proto/auth/pb"

	"github.com/kyamalabs/users/internal/constants"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func startMockGrpcServer(t *testing.T) string {
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	server := grpc.NewServer()

	go func() {
		err := server.Serve(lis)
		require.NoError(t, err)
	}()

	return lis.Addr().String()
}

func TestNewAuthServiceGrpcClient(t *testing.T) {
	testCases := []struct {
		name                 string
		getMockServerAddress func(t *testing.T) string
		checkResponse        func(t *testing.T, client AuthGrpcService, err error)
	}{
		{
			name:                 "successfully creates auth service gRPC client",
			getMockServerAddress: startMockGrpcServer,
			checkResponse: func(t *testing.T, client AuthGrpcService, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, client)
				require.IsType(t, &AuthServiceGrpcClient{}, client)
			},
		},
		{
			name: "connection error",
			getMockServerAddress: func(t *testing.T) string {
				return ""
			},
			checkResponse: func(t *testing.T, client AuthGrpcService, err error) {
				require.Error(t, err)
				require.Nil(t, client)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			privateKeys := []string{"privateKey1", "privateKey2"}

			client, err := NewAuthServiceGrpcClient(tc.getMockServerAddress(t), privateKeys)
			tc.checkResponse(t, client, err)

			if client != nil {
				defer func() {
					err = client.(*AuthServiceGrpcClient).Close()
					require.NoError(t, err)
				}()
			}
		})
	}
}

func TestAuthServiceGrpcClient_withMetadata(t *testing.T) {
	testCases := []struct {
		name                              string
		serviceAuthenticationPayload      string
		serviceAuthenticationPayloadError error
		requestMetadata                   *requestMetadata
		checkResponse                     func(t *testing.T, ctx context.Context, err error)
	}{
		{
			name:                              "could not get service authentication payload",
			serviceAuthenticationPayload:      "",
			serviceAuthenticationPayloadError: errors.New("some error"),
			requestMetadata:                   &requestMetadata{},
			checkResponse: func(t *testing.T, ctx context.Context, err error) {
				require.Error(t, err)
				require.Nil(t, ctx)
			},
		},
		{
			name:                              "with empty request metadata",
			serviceAuthenticationPayload:      "some-authentication-payload",
			serviceAuthenticationPayloadError: nil,
			requestMetadata:                   &requestMetadata{},
			checkResponse: func(t *testing.T, ctx context.Context, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, ctx)

				md, ok := metadata.FromOutgoingContext(ctx)
				require.True(t, ok)

				serviceAuthenticationHeader := md.Get(constants.XServiceAuthenticationHeader)
				require.NotEmpty(t, serviceAuthenticationHeader)

				authorizationHeader := md.Get(constants.AuthorizationHeader)
				require.Empty(t, authorizationHeader)
			},
		},
		{
			name:                              "with populated request metadata",
			serviceAuthenticationPayload:      "some-authentication-payload",
			serviceAuthenticationPayloadError: nil,
			requestMetadata: &requestMetadata{
				accessToken: "some-access-token",
			},
			checkResponse: func(t *testing.T, ctx context.Context, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, ctx)

				md, ok := metadata.FromOutgoingContext(ctx)
				require.True(t, ok)

				serviceAuthenticationHeader := md.Get(constants.XServiceAuthenticationHeader)
				require.NotEmpty(t, serviceAuthenticationHeader)

				authorizationHeader := md.Get(constants.AuthorizationHeader)
				require.NotEmpty(t, authorizationHeader)
				require.Equal(t, []string{fmt.Sprintf("%s %s", constants.AuthorizationBearer, "some-access-token")}, authorizationHeader)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialGetServiceAuthenticationPayload := getServiceAuthenticationPayload
			getServiceAuthenticationPayload = func(serviceAuthPrivateKeys []string) (string, error) {
				return tc.serviceAuthenticationPayload, tc.serviceAuthenticationPayloadError
			}
			defer func() {
				getServiceAuthenticationPayload = initialGetServiceAuthenticationPayload
			}()

			privateKeys := []string{"privateKey1", "privateKey2"}

			client, err := NewAuthServiceGrpcClient(startMockGrpcServer(t), privateKeys)
			require.NoError(t, err)
			require.NotEmpty(t, client)

			defer func() {
				err := client.(*AuthServiceGrpcClient).Close()
				require.NoError(t, err)
			}()

			ctx, err := client.(*AuthServiceGrpcClient).withMetadata(context.Background(), tc.requestMetadata)
			tc.checkResponse(t, ctx, err)
		})
	}
}

type mockAuthClientHandler struct {
	IsVerifyAccessTokenCalled bool
}

func (h *mockAuthClientHandler) GetChallenge(_ context.Context, _ *pb.GetChallengeRequest, _ ...grpc.CallOption) (*pb.GetChallengeResponse, error) {
	return nil, nil
}

func (h *mockAuthClientHandler) AuthenticateAccount(_ context.Context, _ *pb.AuthenticateAccountRequest, _ ...grpc.CallOption) (*pb.AuthenticateAccountResponse, error) {
	return nil, nil
}

func (h *mockAuthClientHandler) VerifyAccessToken(_ context.Context, _ *pb.VerifyAccessTokenRequest, _ ...grpc.CallOption) (*pb.VerifyAccessTokenResponse, error) {
	h.IsVerifyAccessTokenCalled = true
	return nil, nil
}

func (h *mockAuthClientHandler) RefreshAccessToken(_ context.Context, _ *pb.RefreshAccessTokenRequest, _ ...grpc.CallOption) (*pb.RefreshAccessTokenResponse, error) {
	return nil, nil
}

func (h *mockAuthClientHandler) RevokeRefreshTokens(_ context.Context, _ *pb.RevokeRefreshTokensRequest, _ ...grpc.CallOption) (*pb.RevokeRefreshTokensResponse, error) {
	return nil, nil
}

func TestAuthServiceGrpcClient_VerifyAccessToken(t *testing.T) {
	dummyClient := &mockAuthClientHandler{}

	mockAuthServiceGrpcClient := &AuthServiceGrpcClient{
		client:                 dummyClient,
		serviceAuthPrivateKeys: []string{"privateKey1", "privateKey2"},
	}

	initialGetServiceAuthenticationPayload := getServiceAuthenticationPayload
	getServiceAuthenticationPayload = func(serviceAuthPrivateKeys []string) (string, error) {
		return "some authentication payload", nil
	}
	defer func() {
		getServiceAuthenticationPayload = initialGetServiceAuthenticationPayload
	}()

	payload := &pb.VerifyAccessTokenRequest{
		WalletAddress: "some wallet address",
	}

	_, err := mockAuthServiceGrpcClient.VerifyAccessToken(payload, "some access token")
	require.NoError(t, err)

	require.True(t, dummyClient.IsVerifyAccessTokenCalled)
}
