package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/kyamalabs/auth/pkg/util"
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

func generateGetPublicProfileReqParams(t *testing.T) *pb.GetProfileRequest {
	wallet, err := util.NewEthereumWallet()
	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	return &pb.GetProfileRequest{
		WalletAddress: wallet.Address,
	}
}

func TestGetPublicProfileAPI(t *testing.T) {
	getPublicProfileReqParams := generateGetPublicProfileReqParams(t)
	require.NotEmpty(t, getPublicProfileReqParams)

	testCases := []struct {
		name          string
		req           *pb.GetProfileRequest
		buildContext  func(t *testing.T) context.Context
		buildStubs    func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.GetPublicProfileResponse, err error)
	}{
		{
			name: "success",
			req:  getPublicProfileReqParams,
			buildContext: func(t *testing.T) context.Context {
				return handler.NewContextWithBearerToken()
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					GetProfileTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetProfileTxResult{
						EnsName: "mamabear",
						Profile: db.Profile{
							WalletAddress: getPublicProfileReqParams.GetWalletAddress(),
							GamerTag:      "mamabear",
						},
					}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetPublicProfileResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)

				require.Equal(t, getPublicProfileReqParams.GetWalletAddress(), res.GetProfile().GetWalletAddress())
				require.Equal(t, "mamabear", res.GetProfile().GetEnsName())
				require.NotEmpty(t, res.GetProfile().GetGamerTag())
				require.NotZero(t, res.GetProfile().GetCreatedAt())
			},
		},
		{
			name: "invalid request parameters",
			req: &pb.GetProfileRequest{
				WalletAddress: getPublicProfileReqParams.GetWalletAddress()[:len(getPublicProfileReqParams.GetWalletAddress())-1],
			},
			buildContext: func(t *testing.T) context.Context {
				return nil
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
			},
			checkResponse: func(t *testing.T, res *pb.GetPublicProfileResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				expectedFieldViolations := []string{"wallet_address"}
				handler.CheckInvalidRequestParams(t, err, expectedFieldViolations)
			},
		},
		{
			name: "profile does not exist",
			req:  getPublicProfileReqParams,
			buildContext: func(t *testing.T) context.Context {
				return handler.NewContextWithBearerToken()
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					GetProfileTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetProfileTxResult{}, db.RecordNotFoundError)
			},
			checkResponse: func(t *testing.T, res *pb.GetPublicProfileResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				require.ErrorContains(t, err, DoesNotExist)
			},
		},
		{
			name: "db error",
			req:  getPublicProfileReqParams,
			buildContext: func(t *testing.T) context.Context {
				return handler.NewContextWithBearerToken()
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					GetProfileTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetProfileTxResult{}, errors.New("some db error"))
			},
			checkResponse: func(t *testing.T, res *pb.GetPublicProfileResponse, err error) {
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
			res, err := h.GetPublicProfile(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
