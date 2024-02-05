package db

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/kyamalabs/auth/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestGetProfileTx(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	ethereumWallet, err := util.NewEthereumWallet()
	require.NoError(t, err)
	require.NotEmpty(t, ethereumWallet)

	profile, err := testStore.CreateProfile(context.Background(), CreateProfileParams{
		WalletAddress: ethereumWallet.Address,
		GamerTag:      gofakeit.Gamertag(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, profile)

	isAfterCreateCalled := false

	params := GetProfileTxParams{
		WalletAddress: ethereumWallet.Address,
		AfterCreate: func() (string, error) {
			isAfterCreateCalled = true
			return "", nil
		},
	}

	txResult, err := testStore.GetProfileTx(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, txResult)

	require.Equal(t, params.WalletAddress, txResult.Profile.WalletAddress)
	require.Equal(t, profile.GamerTag, txResult.Profile.GamerTag)
	require.WithinDuration(t, profile.CreatedAt, txResult.Profile.CreatedAt, time.Second)

	require.True(t, isAfterCreateCalled)
}
