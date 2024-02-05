package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
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

func generateProfiles(t *testing.T, num int) []db.Profile {
	var profiles []db.Profile

	for i := 0; i < num; i++ {
		wallet, err := util.NewEthereumWallet()
		require.NoError(t, err)
		require.NotEmpty(t, wallet)

		profiles = append(profiles, db.Profile{
			WalletAddress: wallet.Address,
			GamerTag:      gofakeit.Gamertag(),
		})
	}

	return profiles
}

func TestListAllProfilesAPI(t *testing.T) {
	testCases := []struct {
		name          string
		req           *pb.ListAllProfilesRequest
		buildStubs    func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.ListAllProfilesResponse, err error)
	}{
		{
			name: "success",
			req:  &pb.ListAllProfilesRequest{},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					ListProfiles(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generateProfiles(t, 2), nil)

				store.EXPECT().
					GetProfilesCount(gomock.Any()).
					Times(1).
					Return(int64(2), nil)

				cache.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Times(2).
					Return(gofakeit.Gamertag(), nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListAllProfilesResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)

				require.Len(t, res.GetProfiles(), 2)
				require.Equal(t, int32(2), res.GetTotalProfiles())
				require.Equal(t, int32(1), res.GetPage())
				require.Equal(t, int32(30), res.GetPageSize())
			},
		},
		{
			name: "invalid request parameters",
			req: &pb.ListAllProfilesRequest{
				Page:     3,
				PageSize: 55,
			},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
			},
			checkResponse: func(t *testing.T, res *pb.ListAllProfilesResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				expectedFieldViolations := []string{"page_size"}
				checkInvalidRequestParams(t, err, expectedFieldViolations)
			},
		},
		{
			name: "ListProfiles db error",
			req:  &pb.ListAllProfilesRequest{},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					ListProfiles(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("some db error"))
			},
			checkResponse: func(t *testing.T, res *pb.ListAllProfilesResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				require.ErrorContains(t, err, handler.InternalServerError)
			},
		},
		{
			name: "GetProfilesCount db error",
			req:  &pb.ListAllProfilesRequest{},
			buildStubs: func(store *mockdb.MockStore, cache *mockcache.MockCache, authService *mockservices.MockAuthGrpcService, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					ListProfiles(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generateProfiles(t, 2), nil)

				store.EXPECT().
					GetProfilesCount(gomock.Any()).
					Times(1).
					Return(int64(0), errors.New("some db error"))
			},
			checkResponse: func(t *testing.T, res *pb.ListAllProfilesResponse, err error) {
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

			res, err := h.ListAllProfiles(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
