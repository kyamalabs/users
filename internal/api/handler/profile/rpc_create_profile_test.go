package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v5/pgconn"
	authPb "github.com/kyamalabs/auth/api/pb"
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

func generateCreateProfileReqParams(t *testing.T) *pb.CreateProfileRequest {
	wallet, err := util.NewEthereumWallet()
	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	return &pb.CreateProfileRequest{
		WalletAddress: wallet.Address,
		GamerTag:      gofakeit.Gamertag(),
	}
}

func TestCreateProfileAPI(t *testing.T) {
	createProfileReqParams := generateCreateProfileReqParams(t)
	require.NotEmpty(t, createProfileReqParams)

	testCases := []struct {
		name          string
		req           *pb.CreateProfileRequest
		buildContext  func(t *testing.T) context.Context
		buildStubs    func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateProfileResponse, err error)
	}{
		{
			name: "success",
			req:  createProfileReqParams,
			buildContext: func(t *testing.T) context.Context {
				return newContextWithBearerToken()
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				authService.EXPECT().
					VerifyAccessToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&authPb.VerifyAccessTokenResponse{
						Payload: &authPb.AccessTokenPayload{
							Id:            "some-id",
							WalletAddress: createProfileReqParams.WalletAddress,
							Role:          authPb.AccessTokenPayload_GAMER,
						},
					}, nil)

				store.EXPECT().
					CreateProfileTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateProfileTxResult{
						Profile: db.Profile{
							WalletAddress: createProfileReqParams.GetWalletAddress(),
							GamerTag:      createProfileReqParams.GetGamerTag(),
						},
					}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateProfileResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)

				require.Equal(t, createProfileReqParams.GetWalletAddress(), res.Profile.GetWalletAddress())
				require.Equal(t, createProfileReqParams.GetGamerTag(), res.Profile.GetGamerTag())
			},
		},
		{
			name: "invalid request arguments",
			req: &pb.CreateProfileRequest{
				WalletAddress: createProfileReqParams.GetWalletAddress()[:len(createProfileReqParams.GetWalletAddress())-1],
				GamerTag:      createProfileReqParams.GetGamerTag()[:2],
			},
			buildContext: func(t *testing.T) context.Context {
				return nil
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
			},
			checkResponse: func(t *testing.T, res *pb.CreateProfileResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				expectedFieldViolations := []string{"wallet_address", "gamer_tag"}
				checkInvalidRequestParams(t, err, expectedFieldViolations)
			},
		},
		{
			name: "unauthorized user",
			req:  createProfileReqParams,
			buildContext: func(t *testing.T) context.Context {
				return newContextWithBearerToken()
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				authService.EXPECT().
					VerifyAccessToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("some verify access token error"))
			},
			checkResponse: func(t *testing.T, res *pb.CreateProfileResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				require.ErrorContains(t, err, handler.UnauthorizedAccessError)
			},
		},
		{
			name: "profile already exists",
			req:  createProfileReqParams,
			buildContext: func(t *testing.T) context.Context {
				return newContextWithBearerToken()
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				authService.EXPECT().
					VerifyAccessToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&authPb.VerifyAccessTokenResponse{
						Payload: &authPb.AccessTokenPayload{
							Id:            "some-id",
							WalletAddress: createProfileReqParams.WalletAddress,
							Role:          authPb.AccessTokenPayload_GAMER,
						},
					}, nil)

				store.EXPECT().
					CreateProfileTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateProfileTxResult{}, &pgconn.PgError{Code: db.UniqueViolation})
			},
			checkResponse: func(t *testing.T, res *pb.CreateProfileResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				require.ErrorContains(t, err, AlreadyExists)
			},
		},
		{
			name: "db error",
			req:  createProfileReqParams,
			buildContext: func(t *testing.T) context.Context {
				return newContextWithBearerToken()
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				authService.EXPECT().
					VerifyAccessToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&authPb.VerifyAccessTokenResponse{
						Payload: &authPb.AccessTokenPayload{
							Id:            "some-id",
							WalletAddress: createProfileReqParams.WalletAddress,
							Role:          authPb.AccessTokenPayload_GAMER,
						},
					}, nil)

				store.EXPECT().
					CreateProfileTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateProfileTxResult{}, errors.New("some db error"))
			},
			checkResponse: func(t *testing.T, res *pb.CreateProfileResponse, err error) {
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
			res, err := h.CreateProfile(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
