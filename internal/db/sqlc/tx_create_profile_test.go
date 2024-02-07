package db

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/kyamalabs/auth/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestCreateProfileTx(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	ethereumWallet, err := util.NewEthereumWallet()
	require.NoError(t, err)
	require.NotEmpty(t, ethereumWallet)

	isAfterCreateCalled := false

	params := CreateProfileTxParams{
		AfterCreate: func() error {
			isAfterCreateCalled = true
			return nil
		},
		CreateProfileParams: CreateProfileParams{
			WalletAddress: ethereumWallet.Address,
			GamerTag:      gofakeit.Gamertag(),
		},
	}

	txResult, err := testStore.CreateProfileTx(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, txResult)

	require.Equal(t, params.WalletAddress, txResult.Profile.WalletAddress)
	require.Equal(t, params.GamerTag, txResult.Profile.GamerTag)
	require.NotZero(t, txResult.Profile.CreatedAt)

	require.True(t, isAfterCreateCalled)
}
