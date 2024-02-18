package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/kyamalabs/auth/pkg/util"
	authPb "github.com/kyamalabs/proto/proto/auth/pb"
	"github.com/kyamalabs/users/api/pb"
	"github.com/kyamalabs/users/internal/api/handler"
	mockcache "github.com/kyamalabs/users/internal/cache/mock"
	mockdb "github.com/kyamalabs/users/internal/db/mock"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	mockservices "github.com/kyamalabs/users/internal/services/mock"
	mockwk "github.com/kyamalabs/users/internal/worker/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func generateUpdateProfileReqParams(t *testing.T) *pb.UpdateProfileRequest {
	wallet, err := util.NewEthereumWallet()
	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	return &pb.UpdateProfileRequest{
		WalletAddress: wallet.Address,
		GamerTag:      gofakeit.Gamertag(),
	}
}

func TestUpdateProfileAPI(t *testing.T) {
	updateProfileReqParams := generateUpdateProfileReqParams(t)
	require.NotEmpty(t, updateProfileReqParams)

	testCases := []struct {
		name          string
		req           *pb.UpdateProfileRequest
		buildContext  func(t *testing.T) context.Context
		buildStubs    func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.UpdateProfileResponse, err error)
	}{
		{
			name: "success",
			req:  updateProfileReqParams,
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
							WalletAddress: updateProfileReqParams.WalletAddress,
							Role:          authPb.AccessTokenPayload_GAMER,
						},
					}, nil)

				store.EXPECT().
					UpdateProfileTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.UpdateProfileTxResult{
						EnsName: "mamabear",
						Profile: db.Profile{
							WalletAddress: updateProfileReqParams.GetWalletAddress(),
							GamerTag:      updateProfileReqParams.GetGamerTag(),
						},
					}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateProfileResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)

				require.Equal(t, updateProfileReqParams.GetWalletAddress(), res.GetProfile().GetWalletAddress())
				require.Equal(t, "mamabear", res.GetProfile().GetEnsName())
				require.Equal(t, updateProfileReqParams.GamerTag, res.GetProfile().GetGamerTag())
				require.NotZero(t, res.GetProfile().GetCreatedAt())
			},
		},
		{
			name: "invalid request parameters",
			req: &pb.UpdateProfileRequest{
				WalletAddress: updateProfileReqParams.GetWalletAddress()[:len(updateProfileReqParams.GetWalletAddress())-1],
				GamerTag:      updateProfileReqParams.GetGamerTag()[:2],
			},
			buildContext: func(t *testing.T) context.Context {
				return nil
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
			},
			checkResponse: func(t *testing.T, res *pb.UpdateProfileResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				expectedFieldViolations := []string{"wallet_address", "gamer_tag"}
				handler.CheckInvalidRequestParams(t, err, expectedFieldViolations)
			},
		},
		{
			name: "unauthorized user",
			req:  updateProfileReqParams,
			buildContext: func(t *testing.T) context.Context {
				return handler.NewContextWithBearerToken()
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				authService.EXPECT().
					VerifyAccessToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("some verify access token error"))
			},
			checkResponse: func(t *testing.T, res *pb.UpdateProfileResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				require.ErrorContains(t, err, handler.UnauthorizedAccessError)
			},
		},
		{
			name: "db error",
			req:  updateProfileReqParams,
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					UpdateProfileTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.UpdateProfileTxResult{}, errors.New("some db error"))

				authService.EXPECT().
					VerifyAccessToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&authPb.VerifyAccessTokenResponse{
						Payload: &authPb.AccessTokenPayload{
							Id:            "some-id",
							WalletAddress: updateProfileReqParams.WalletAddress,
							Role:          authPb.AccessTokenPayload_GAMER,
						},
					}, nil)
			},
			buildContext: func(t *testing.T) context.Context {
				return handler.NewContextWithBearerToken()
			},
			checkResponse: func(t *testing.T, res *pb.UpdateProfileResponse, err error) {
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
			res, err := h.UpdateProfile(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
