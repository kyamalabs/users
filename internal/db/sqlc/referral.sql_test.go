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

func TestGetReferralsCount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	referrer := createTestProfile(t)

	initialReferralsCount, err := testStore.GetReferralsCount(context.Background(), referrer.WalletAddress)
	require.NoError(t, err)
	require.NotNil(t, initialReferralsCount)

	numAdditionalReferrals := 4
	for i := 0; i < numAdditionalReferrals; i++ {
		profile := createTestProfile(t)
		require.NotEmpty(t, profile)

		referral, err := testStore.CreateReferral(context.Background(), CreateReferralParams{
			Referrer: referrer.WalletAddress,
			Referee:  profile.WalletAddress,
		})

		require.NotEmpty(t, referral)
		require.NoError(t, err)
	}

	finalReferralsCount, err := testStore.GetReferralsCount(context.Background(), referrer.WalletAddress)
	require.NoError(t, err)
	require.NotEmpty(t, finalReferralsCount)

	require.Equal(t, initialReferralsCount+int64(numAdditionalReferrals), finalReferralsCount)
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
