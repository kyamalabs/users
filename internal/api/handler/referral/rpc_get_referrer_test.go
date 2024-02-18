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

func generateGetReferrerReqParams(t *testing.T) *pb.GetReferrerRequest {
	wallet, err := util.NewEthereumWallet()
	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	return &pb.GetReferrerRequest{
		WalletAddress: wallet.Address,
	}
}

func TestGetReferrerAPI(t *testing.T) {
	getReferrerReqParams := generateGetReferrerReqParams(t)
	require.NotEmpty(t, getReferrerReqParams)

	testCases := []struct {
		name          string
		req           *pb.GetReferrerRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.GetReferrerResponse, err error)
	}{
		{
			name: "success",
			req:  getReferrerReqParams,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetReferrer(gomock.Any(), gomock.Eq(getReferrerReqParams.GetWalletAddress())).
					Times(1).
					Return(db.Referral{
						Referrer: "some-wallet-address",
						Referee:  getReferrerReqParams.GetWalletAddress(),
					}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetReferrerResponse, err error) {
				require.NoError(t, err)

				require.Equal(t, "some-wallet-address", res.GetReferral().GetReferrer())
				require.Equal(t, getReferrerReqParams.GetWalletAddress(), res.GetReferral().GetReferee())
			},
		},
		{
			name: "invalid request parameters",
			req: &pb.GetReferrerRequest{
				WalletAddress: getReferrerReqParams.GetWalletAddress()[:len(getReferrerReqParams.GetWalletAddress())-1],
			},
			buildStubs: func(store *mockdb.MockStore) {
			},
			checkResponse: func(t *testing.T, res *pb.GetReferrerResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				expectedFieldViolations := []string{"wallet_address"}
				handler.CheckInvalidRequestParams(t, err, expectedFieldViolations)
			},
		},
		{
			name: "db error",
			req:  getReferrerReqParams,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetReferrer(gomock.Any(), gomock.Eq(getReferrerReqParams.GetWalletAddress())).
					Times(1).
					Return(db.Referral{}, errors.New("some db error"))
			},
			checkResponse: func(t *testing.T, res *pb.GetReferrerResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)

				require.ErrorContains(t, err, handler.InternalServerError)
			},
		},
		{
			name: "referral not present in db",
			req:  getReferrerReqParams,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetReferrer(gomock.Any(), gomock.Eq(getReferrerReqParams.GetWalletAddress())).
					Times(1).
					Return(db.Referral{}, db.RecordNotFoundError)
			},
			checkResponse: func(t *testing.T, res *pb.GetReferrerResponse, err error) {
				require.NoError(t, err)

				require.Equal(t, "", res.GetReferral().GetReferrer())
				require.Equal(t, "", res.GetReferral().GetReferee())
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

			res, err := h.GetReferrer(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
