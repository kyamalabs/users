package middleware

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/kyamalabs/auth/api/pb"
	"github.com/kyamalabs/users/internal/constants"
	mockservices "github.com/kyamalabs/users/internal/services/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

func TestAuthorizeUser(t *testing.T) {
	testCases := []struct {
		name          string
		buildStubs    func(authService *mockservices.MockAuthGrpcService)
		buildContext  func(t *testing.T) context.Context
		walletAddress string
		checkResponse func(t *testing.T, payload *pb.AccessTokenPayload, err error)
	}{
		{
			name: "success",
			buildStubs: func(authService *mockservices.MockAuthGrpcService) {
				expectedPayload := &pb.VerifyAccessTokenRequest{
					WalletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
				}

				response := &pb.VerifyAccessTokenResponse{
					Payload: &pb.AccessTokenPayload{
						Id:            "some-id",
						WalletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
						Role:          pb.AccessTokenPayload_GAMER,
					},
				}

				authService.EXPECT().
					VerifyAccessToken(expectedPayload, "some-dummy-access-token").
					Times(1).
					Return(response, nil)
			},
			buildContext: func(t *testing.T) context.Context {
				md := metadata.MD{
					constants.AuthorizationHeader: []string{
						fmt.Sprintf("%s %s", constants.AuthorizationBearer, "some-dummy-access-token"),
					},
				}

				return metadata.NewIncomingContext(context.Background(), md)
			},
			walletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
			checkResponse: func(t *testing.T, payload *pb.AccessTokenPayload, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, payload)

				require.Equal(t, payload.Id, "some-id")
				require.Equal(t, payload.WalletAddress, "0xc0ffee254729296a45a3885639AC7E10F9d54979")
				require.Equal(t, payload.Role, pb.AccessTokenPayload_GAMER)
			},
		},
		{
			name: "missing metadata from incoming context",
			buildStubs: func(authService *mockservices.MockAuthGrpcService) {
			},
			buildContext: func(t *testing.T) context.Context {
				return context.Background()
			},
			walletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
			checkResponse: func(t *testing.T, payload *pb.AccessTokenPayload, err error) {
				require.Error(t, err)
				require.Nil(t, payload)
			},
		},
		{
			name: "missing authorization header",
			buildStubs: func(authService *mockservices.MockAuthGrpcService) {
			},
			buildContext: func(t *testing.T) context.Context {
				md := metadata.MD{
					"some_other_header": []string{
						"some_value",
					},
				}

				return metadata.NewIncomingContext(context.Background(), md)
			},
			walletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
			checkResponse: func(t *testing.T, payload *pb.AccessTokenPayload, err error) {
				require.Error(t, err)
				require.Nil(t, payload)
			},
		},
		{
			name: "invalid authorization header format",
			buildStubs: func(authService *mockservices.MockAuthGrpcService) {
			},
			buildContext: func(t *testing.T) context.Context {
				md := metadata.MD{
					constants.AuthorizationHeader: []string{
						"some_value",
					},
				}

				return metadata.NewIncomingContext(context.Background(), md)
			},
			walletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
			checkResponse: func(t *testing.T, payload *pb.AccessTokenPayload, err error) {
				require.Error(t, err)
				require.Nil(t, payload)
			},
		},
		{
			name: "unsupported authorization type",
			buildStubs: func(authService *mockservices.MockAuthGrpcService) {
			},
			buildContext: func(t *testing.T) context.Context {
				md := metadata.MD{
					constants.AuthorizationHeader: []string{
						fmt.Sprintf("%s %s", "unsupported_auth_type", "some_token"),
					},
				}

				return metadata.NewIncomingContext(context.Background(), md)
			},
			walletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
			checkResponse: func(t *testing.T, payload *pb.AccessTokenPayload, err error) {
				require.Error(t, err)
				require.Nil(t, payload)
			},
		},
		{
			name: "invalid access token",
			buildStubs: func(authService *mockservices.MockAuthGrpcService) {
				expectedPayload := &pb.VerifyAccessTokenRequest{
					WalletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
				}

				authService.EXPECT().
					VerifyAccessToken(expectedPayload, "some-dummy-access-token").
					Times(1).
					Return(nil, errors.New("invalid token error"))
			},
			buildContext: func(t *testing.T) context.Context {
				md := metadata.MD{
					constants.AuthorizationHeader: []string{
						fmt.Sprintf("%s %s", constants.AuthorizationBearer, "some-dummy-access-token"),
					},
				}

				return metadata.NewIncomingContext(context.Background(), md)
			},
			walletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
			checkResponse: func(t *testing.T, payload *pb.AccessTokenPayload, err error) {
				require.Error(t, err)
				require.Nil(t, payload)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			authService := mockservices.NewMockAuthGrpcService(ctrl)

			tc.buildStubs(authService)

			ctx := tc.buildContext(t)
			payload, err := AuthorizeUser(ctx, tc.walletAddress, authService)
			tc.checkResponse(t, payload, err)
		})
	}
}
