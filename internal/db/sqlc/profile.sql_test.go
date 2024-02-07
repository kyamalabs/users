package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/kyamalabs/auth/pkg/util"
	"github.com/stretchr/testify/require"
)

func createTestProfile(t *testing.T) Profile {
	ethereumWallet, err := util.NewEthereumWallet()
	require.NoError(t, err)
	require.NotEmpty(t, ethereumWallet)

	params := CreateProfileParams{
		WalletAddress: ethereumWallet.Address,
		GamerTag:      gofakeit.Gamertag(),
	}

	profile, err := testStore.CreateProfile(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, profile)

	require.Equal(t, params.WalletAddress, profile.WalletAddress)
	require.Equal(t, params.GamerTag, profile.GamerTag)
	require.NotZero(t, profile.CreatedAt)

	return profile
}

func TestCreateProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	profile := createTestProfile(t)
	require.NotEmpty(t, profile)
}

func TestGetProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	profile := createTestProfile(t)
	require.NotEmpty(t, profile)

	fetchedProfile, err := testStore.GetProfile(context.Background(), profile.WalletAddress)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedProfile)

	require.Equal(t, profile.WalletAddress, fetchedProfile.WalletAddress)
	require.Equal(t, profile.GamerTag, fetchedProfile.GamerTag)
	require.WithinDuration(t, profile.CreatedAt, fetchedProfile.CreatedAt, time.Second)
}

func TestGetProfilesCount(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	initialProfilesCount, err := testStore.GetProfilesCount(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, initialProfilesCount)

	numAdditionalProfiles := 6
	for i := 0; i < numAdditionalProfiles; i++ {
		profile := createTestProfile(t)
		require.NotEmpty(t, profile)
	}

	finalProfilesCount, err := testStore.GetProfilesCount(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, finalProfilesCount)

	require.Equal(t, initialProfilesCount+int64(numAdditionalProfiles), finalProfilesCount)
	fmt.Println(finalProfilesCount)
}

func TestListProfiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	var recentWalletAddresses []string
	numProfilesToCreate := 4

	for i := 0; i < numProfilesToCreate; i++ {
		profile := createTestProfile(t)
		require.NotEmpty(t, profile)
		recentWalletAddresses = append(recentWalletAddresses, profile.WalletAddress)
	}

	params := ListProfilesParams{
		Limit:  int32(numProfilesToCreate - 1),
		Offset: 1,
	}
	recentlyCreatedProfiles, err := testStore.ListProfiles(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, recentlyCreatedProfiles)

	for idx, profile := range recentlyCreatedProfiles {
		require.Equal(t, recentWalletAddresses[len(recentWalletAddresses)-idx-2], profile.WalletAddress)
	}
}

func TestDeleteProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	profile := createTestProfile(t)
	require.NotEmpty(t, profile)

	err := testStore.DeleteProfile(context.Background(), profile.WalletAddress)
	require.NoError(t, err)

	deletedProfile, err := testStore.GetProfile(context.Background(), profile.WalletAddress)
	require.Error(t, err)
	require.Empty(t, deletedProfile)

	require.Equal(t, RecordNotFoundError, err)
}

func TestUpdateProfile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test to maintain db state")
	}

	profile := createTestProfile(t)
	require.NotEmpty(t, profile)

	expectedUpdatedGamerTag := fmt.Sprintf("%s--updated", profile.GamerTag)

	params := UpdateProfileParams{
		WalletAddress: profile.WalletAddress,
		GamerTag: pgtype.Text{
			String: expectedUpdatedGamerTag,
			Valid:  true,
		},
	}
	updatedProfile, err := testStore.UpdateProfile(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, updatedProfile)

	require.Equal(t, expectedUpdatedGamerTag, updatedProfile.GamerTag)
	require.Equal(t, profile.WalletAddress, updatedProfile.WalletAddress)
	require.WithinDuration(t, profile.CreatedAt, updatedProfile.CreatedAt, time.Second)
}
