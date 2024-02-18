package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/kyamalabs/auth/pkg/util"
	authPb "github.com/kyamalabs/proto/proto/auth/pb"
	"github.com/kyamalabs/users/api/pb"
	"github.com/kyamalabs/users/internal/api/handler"
	mockcache "github.com/kyamalabs/users/internal/cache/mock"
	mockdb "github.com/kyamalabs/users/internal/db/mock"
	mockservices "github.com/kyamalabs/users/internal/services/mock"
	mockwk "github.com/kyamalabs/users/internal/worker/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
)

func generateDeleteProfileReqParams(t *testing.T) *pb.DeleteProfileRequest {
	wallet, err := util.NewEthereumWallet()
	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	return &pb.DeleteProfileRequest{
		WalletAddress: wallet.Address,
	}
}

func TestDeleteProfileAPI(t *testing.T) {
	deleteProfileReqParams := generateDeleteProfileReqParams(t)
	require.NotEmpty(t, deleteProfileReqParams)

	testCases := []struct {
		name          string
		req           *pb.DeleteProfileRequest
		buildContext  func(t *testing.T) context.Context
		buildStubs    func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *emptypb.Empty, err error)
	}{
		{
			name: "success",
			req:  deleteProfileReqParams,
			buildContext: func(t *testing.T) context.Context {
				return handler.NewContextWithBearerToken()
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				authService.EXPECT().
					VerifyAccessToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&authPb.VerifyAccessTokenResponse{
						Payload: &authPb.AccessTokenPayload{
							Id:            "some-id",
							WalletAddress: deleteProfileReqParams.WalletAddress,
							Role:          authPb.AccessTokenPayload_GAMER,
						},
					}, nil)

				store.EXPECT().
					DeleteProfile(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *emptypb.Empty, err error) {
				require.NoError(t, err)
				require.Equal(t, &emptypb.Empty{}, res)
			},
		},
		{
			name: "invalid request arguments",
			req: &pb.DeleteProfileRequest{
				WalletAddress: deleteProfileReqParams.GetWalletAddress()[:len(deleteProfileReqParams.GetWalletAddress())-1],
			},
			buildContext: func(t *testing.T) context.Context {
				return nil
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
			},
			checkResponse: func(t *testing.T, res *emptypb.Empty, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				expectedFieldViolations := []string{"wallet_address"}
				handler.CheckInvalidRequestParams(t, err, expectedFieldViolations)
			},
		},
		{
			name: "unauthorized user",
			req:  deleteProfileReqParams,
			buildContext: func(t *testing.T) context.Context {
				return handler.NewContextWithBearerToken()
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				authService.EXPECT().
					VerifyAccessToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("some verify access token error"))
			},
			checkResponse: func(t *testing.T, res *emptypb.Empty, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				require.ErrorContains(t, err, handler.UnauthorizedAccessError)
			},
		},
		{
			name: "db error",
			req:  deleteProfileReqParams,
			buildContext: func(t *testing.T) context.Context {
				return handler.NewContextWithBearerToken()
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				authService.EXPECT().
					VerifyAccessToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&authPb.VerifyAccessTokenResponse{
						Payload: &authPb.AccessTokenPayload{
							Id:            "some-id",
							WalletAddress: deleteProfileReqParams.WalletAddress,
							Role:          authPb.AccessTokenPayload_GAMER,
						},
					}, nil)

				store.EXPECT().
					DeleteProfile(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("some db error"))
			},
			checkResponse: func(t *testing.T, res *emptypb.Empty, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				require.ErrorContains(t, err, handler.InternalServerError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			cache := mockcache.NewMockCache(ctrl)
			authService := mockservices.NewMockAuthGrpcService(ctrl)
			taskDistributor := mockwk.NewMockTaskDistributor(ctrl)

			tc.buildStubs(store, cache, authService, taskDistributor)

			h := newTestHandler(store, cache, authService, taskDistributor)

			ctx := tc.buildContext(t)
			res, err := h.DeleteProfile(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
