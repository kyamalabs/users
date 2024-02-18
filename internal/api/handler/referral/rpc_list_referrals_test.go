package referral

import (
	"context"
	"errors"
	"testing"

	"github.com/kyamalabs/auth/pkg/util"
	"github.com/kyamalabs/users/api/pb"
	"github.com/kyamalabs/users/internal/api/handler"
	mockdb "github.com/kyamalabs/users/internal/db/mock"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func generateListReferralsReqParams(t *testing.T) *pb.ListReferralsRequest {
	wallet, err := util.NewEthereumWallet()
	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	return &pb.ListReferralsRequest{
		WalletAddress: wallet.Address,
	}
}

func generateReferrals(t *testing.T, num int) []db.Referral {
	var referrals []db.Referral

	for i := 0; i < num; i++ {
		referrer, err := util.NewEthereumWallet()
		require.NoError(t, err)
		require.NotEmpty(t, referrer)

		referee, err := util.NewEthereumWallet()
		require.NoError(t, err)
		require.NotEmpty(t, referee)

		referrals = append(referrals, db.Referral{
			Referrer: referrer.Address,
			Referee:  referee.Address,
		})
	}

	return referrals
}

func TestListReferralsAPI(t *testing.T) {
	listReferralsReqParams := generateListReferralsReqParams(t)
	require.NotEmpty(t, listReferralsReqParams)

	testCases := []struct {
		name          string
		req           *pb.ListReferralsRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.ListReferralsResponse, err error)
	}{
		{
			name: "success",
			req:  listReferralsReqParams,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListReferrals(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generateReferrals(t, 3), nil)

				store.EXPECT().
					GetReferralsCount(gomock.Any(), gomock.Eq(listReferralsReqParams.GetWalletAddress())).
					Times(1).
					Return(int64(3), nil)
			},
			checkResponse: func(t *testing.T, res *pb.ListReferralsResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)

				require.Len(t, res.GetReferrals(), 3)
				require.Equal(t, int32(3), res.GetTotalReferrals())
				require.Equal(t, int32(1), res.GetPage())
				require.Equal(t, int32(30), res.GetPageSize())
			},
		},
		{
			name: "invalid request parameters",
			req: &pb.ListReferralsRequest{
				Page:          1,
				PageSize:      55,
				WalletAddress: listReferralsReqParams.GetWalletAddress()[:len(listReferralsReqParams.GetWalletAddress())-1],
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, res *pb.ListReferralsResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				expectedFieldViolations := []string{"wallet_address", "page_size"}
				handler.CheckInvalidRequestParams(t, err, expectedFieldViolations)
			},
		},
		{
			name: "ListReferrals db error",
			req:  listReferralsReqParams,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListReferrals(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("some db error"))
			},
			checkResponse: func(t *testing.T, res *pb.ListReferralsResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				require.ErrorContains(t, err, handler.InternalServerError)
			},
		},
		{
			name: "GetReferralsCount db error",
			req:  listReferralsReqParams,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListReferrals(gomock.Any(), gomock.Any()).
					Times(1).
					Return(generateReferrals(t, 3), nil)

				store.EXPECT().
					GetReferralsCount(gomock.Any(), gomock.Eq(listReferralsReqParams.GetWalletAddress())).
					Times(1).
					Return(int64(0), errors.New("some db error"))
			},
			checkResponse: func(t *testing.T, res *pb.ListReferralsResponse, err error) {
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

			tc.buildStubs(store)

			h := newTestHandler(store)

			res, err := h.ListReferrals(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
