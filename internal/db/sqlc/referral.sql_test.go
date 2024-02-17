package db

import (
	"context"
	"testing"
	"time"

	"github.com/kyamalabs/auth/pkg/util"
	"github.com/stretchr/testify/require"
)

func createTestReferral(t *testing.T, referrer Profile) Referral {
	refereeEthereumWallet, err := util.NewEthereumWallet()
	require.NoError(t, err)
	require.NotEmpty(t, refereeEthereumWallet)

	params := CreateReferralParams{
		Referrer: referrer.WalletAddress,
		Referee:  refereeEthereumWallet.Address,
	}

	referral, err := testStore.CreateReferral(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, referral)

	require.NotEmpty(t, referral.ID)
	require.Equal(t, referral.Referrer, referrer.WalletAddress)
	require.Equal(t, referral.Referee, refereeEthereumWallet.Address)
	require.NotZero(t, referral.ReferredAt)

	return referral
}

func TestCreateReferral(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	createTestReferral(t, createTestProfile(t))
}

func TestGetReferrer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	referral := createTestReferral(t, createTestProfile(t))

	fetchedReferral, err := testStore.GetReferrer(context.Background(), referral.Referee)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedReferral)

	require.Equal(t, referral.ID, fetchedReferral.ID)
	require.Equal(t, referral.Referrer, fetchedReferral.Referrer)
	require.Equal(t, referral.Referee, fetchedReferral.Referee)
	require.WithinDuration(t, referral.ReferredAt, fetchedReferral.ReferredAt, time.Second)
}

func TestListReferrals(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	referrer := createTestProfile(t)
	numReferralsToCreate := 4

	for i := 0; i < numReferralsToCreate; i++ {
		createTestReferral(t, referrer)
	}

	params := ListReferralsParams{
		Referrer: referrer.WalletAddress,
		Limit:    int32(numReferralsToCreate),
		Offset:   0,
	}

	referrals, err := testStore.ListReferrals(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, referrals)

	require.Len(t, referrals, numReferralsToCreate)
	for _, referral := range referrals {
		require.Equal(t, referrer.WalletAddress, referral.Referrer)
	}
}
