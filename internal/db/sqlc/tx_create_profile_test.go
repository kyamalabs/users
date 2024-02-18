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

	testCases := []struct {
		name        string
		buildParams func(
			t *testing.T,
			newProfileEthereumWallet *util.EthereumWallet,
			referrerEthereumWallet *util.EthereumWallet,
			newProfileGamerTag string,
			isAfterCreateCalled *bool) CreateProfileTxParams
		checkResponse func(
			t *testing.T,
			newProfileEthereumWallet *util.EthereumWallet,
			referrerEthereumWallet *util.EthereumWallet,
			newProfileGamerTag string,
			isAfterCreateCalled *bool,
			txResult CreateProfileTxResult,
			err error)
	}{
		{
			name: "create profile without referer",
			buildParams: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string,
				isAfterCreateCalled *bool,
			) CreateProfileTxParams {
				return CreateProfileTxParams{
					AfterCreate: func() error {
						*isAfterCreateCalled = true
						return nil
					},
					CreateProfileParams: CreateProfileParams{
						WalletAddress: newProfileEthereumWallet.Address,
						GamerTag:      newProfileGamerTag,
					},
				}
			},
			checkResponse: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string, isAfterCreateCalled *bool,
				txResult CreateProfileTxResult,
				err error,
			) {
				require.Equal(t, newProfileEthereumWallet.Address, txResult.Profile.WalletAddress)
				require.Equal(t, newProfileGamerTag, txResult.Profile.GamerTag)
				require.NotZero(t, txResult.Profile.CreatedAt)

				require.Empty(t, txResult.Referral)

				require.NoError(t, err)

				require.True(t, *isAfterCreateCalled)
			},
		},
		{
			name: "create profile with referrer",
			buildParams: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string,
				isAfterCreateCalled *bool,
			) CreateProfileTxParams {
				referrer, err := testStore.CreateProfile(context.Background(), CreateProfileParams{
					WalletAddress: referrerEthereumWallet.Address,
					GamerTag:      gofakeit.Gamertag(),
				})
				require.NotEmpty(t, referrer)
				require.NoError(t, err)

				return CreateProfileTxParams{
					AfterCreate: func() error {
						*isAfterCreateCalled = true
						return nil
					},
					CreateProfileParams: CreateProfileParams{
						WalletAddress: newProfileEthereumWallet.Address,
						GamerTag:      newProfileGamerTag,
					},
					Referrer: referrerEthereumWallet.Address,
				}
			},
			checkResponse: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string, isAfterCreateCalled *bool,
				txResult CreateProfileTxResult,
				err error,
			) {
				require.Equal(t, newProfileEthereumWallet.Address, txResult.Profile.WalletAddress)
				require.Equal(t, newProfileGamerTag, txResult.Profile.GamerTag)
				require.NotZero(t, txResult.Profile.CreatedAt)

				require.Equal(t, referrerEthereumWallet.Address, txResult.Referral.Referrer)
				require.Equal(t, newProfileEthereumWallet.Address, txResult.Referral.Referee)
				require.NotZero(t, txResult.Referral.ReferredAt)

				require.NoError(t, err)

				require.True(t, *isAfterCreateCalled)
			},
		},
		{
			name: "user profile already exists",
			buildParams: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string,
				isAfterCreateCalled *bool,
			) CreateProfileTxParams {
				profile, err := testStore.CreateProfile(context.Background(), CreateProfileParams{
					WalletAddress: newProfileEthereumWallet.Address,
					GamerTag:      gofakeit.Gamertag(),
				})
				require.NotEmpty(t, profile)
				require.NoError(t, err)

				return CreateProfileTxParams{
					AfterCreate: func() error {
						*isAfterCreateCalled = true
						return nil
					},
					CreateProfileParams: CreateProfileParams{
						WalletAddress: newProfileEthereumWallet.Address,
						GamerTag:      newProfileGamerTag,
					},
				}
			},
			checkResponse: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string, isAfterCreateCalled *bool,
				txResult CreateProfileTxResult,
				err error,
			) {
				require.Empty(t, txResult.Profile)
				require.Empty(t, txResult.Referral)

				require.Equal(t, UserProfileAlreadyExistsError, err)

				require.False(t, *isAfterCreateCalled)
			},
		},
		{
			name: "gamer tag already in use",
			buildParams: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string,
				isAfterCreateCalled *bool,
			) CreateProfileTxParams {
				someEthereumWallet, err := util.NewEthereumWallet()
				require.NoError(t, err)
				require.NotEmpty(t, someEthereumWallet)

				profile, err := testStore.CreateProfile(context.Background(), CreateProfileParams{
					WalletAddress: someEthereumWallet.Address,
					GamerTag:      newProfileGamerTag,
				})
				require.NotEmpty(t, profile)
				require.NoError(t, err)

				return CreateProfileTxParams{
					AfterCreate: func() error {
						*isAfterCreateCalled = true
						return nil
					},
					CreateProfileParams: CreateProfileParams{
						WalletAddress: newProfileEthereumWallet.Address,
						GamerTag:      newProfileGamerTag,
					},
				}
			},
			checkResponse: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string, isAfterCreateCalled *bool,
				txResult CreateProfileTxResult,
				err error,
			) {
				require.Empty(t, txResult.Profile)
				require.Empty(t, txResult.Referral)

				require.Equal(t, GamerTagAlreadyInUseError, err)

				require.False(t, *isAfterCreateCalled)
			},
		},
		{
			name: "self referral",
			buildParams: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string,
				isAfterCreateCalled *bool,
			) CreateProfileTxParams {
				return CreateProfileTxParams{
					AfterCreate: func() error {
						*isAfterCreateCalled = true
						return nil
					},
					CreateProfileParams: CreateProfileParams{
						WalletAddress: newProfileEthereumWallet.Address,
						GamerTag:      newProfileGamerTag,
					},
					Referrer: newProfileEthereumWallet.Address,
				}
			},
			checkResponse: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string, isAfterCreateCalled *bool,
				txResult CreateProfileTxResult,
				err error,
			) {
				require.Empty(t, txResult.Profile)
				require.Empty(t, txResult.Referral)

				require.Equal(t, SelfReferralError, err)

				require.False(t, *isAfterCreateCalled)
			},
		},
		{
			name: "referrer does not exist",
			buildParams: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string,
				isAfterCreateCalled *bool,
			) CreateProfileTxParams {
				return CreateProfileTxParams{
					AfterCreate: func() error {
						*isAfterCreateCalled = true
						return nil
					},
					CreateProfileParams: CreateProfileParams{
						WalletAddress: newProfileEthereumWallet.Address,
						GamerTag:      newProfileGamerTag,
					},
					Referrer: referrerEthereumWallet.Address,
				}
			},
			checkResponse: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string, isAfterCreateCalled *bool,
				txResult CreateProfileTxResult,
				err error,
			) {
				require.Empty(t, txResult.Profile)
				require.Empty(t, txResult.Referral)

				require.Equal(t, ReferrerDoesNotExistError, err)

				require.False(t, *isAfterCreateCalled)
			},
		},
		{
			name: "user already referred",
			buildParams: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string,
				isAfterCreateCalled *bool,
			) CreateProfileTxParams {
				referrer, err := testStore.CreateProfile(context.Background(), CreateProfileParams{
					WalletAddress: referrerEthereumWallet.Address,
					GamerTag:      gofakeit.Gamertag(),
				})
				require.NotEmpty(t, referrer)
				require.NoError(t, err)

				someEthereumWallet, err := util.NewEthereumWallet()
				require.NoError(t, err)
				require.NotEmpty(t, someEthereumWallet)

				someProfile, err := testStore.CreateProfile(context.Background(), CreateProfileParams{
					WalletAddress: someEthereumWallet.Address,
					GamerTag:      gofakeit.Gamertag(),
				})
				require.NotEmpty(t, someProfile)
				require.NoError(t, err)

				referral, err := testStore.CreateReferral(context.Background(), CreateReferralParams{
					Referrer: someProfile.WalletAddress,
					Referee:  newProfileEthereumWallet.Address,
				})
				require.NotEmpty(t, referral)
				require.NoError(t, err)

				return CreateProfileTxParams{
					AfterCreate: func() error {
						*isAfterCreateCalled = true
						return nil
					},
					CreateProfileParams: CreateProfileParams{
						WalletAddress: newProfileEthereumWallet.Address,
						GamerTag:      newProfileGamerTag,
					},
					Referrer: referrerEthereumWallet.Address,
				}
			},
			checkResponse: func(
				t *testing.T,
				newProfileEthereumWallet *util.EthereumWallet,
				referrerEthereumWallet *util.EthereumWallet,
				newProfileGamerTag string, isAfterCreateCalled *bool,
				txResult CreateProfileTxResult,
				err error,
			) {
				require.Empty(t, txResult.Profile)
				require.Empty(t, txResult.Referral)

				require.Equal(t, UserAlreadyReferredError, err)

				require.False(t, *isAfterCreateCalled)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isAfterCreateCalled := false

			newProfileEthereumWallet, err := util.NewEthereumWallet()
			require.NoError(t, err)
			require.NotEmpty(t, newProfileEthereumWallet)

			newProfileGamerTag := gofakeit.Gamertag()

			referrerEthereumWallet, err := util.NewEthereumWallet()
			require.NoError(t, err)
			require.NotEmpty(t, referrerEthereumWallet)

			params := tc.buildParams(t, newProfileEthereumWallet, referrerEthereumWallet, newProfileGamerTag, &isAfterCreateCalled)
			txResult, err := testStore.CreateProfileTx(context.Background(), params)
			tc.checkResponse(t, newProfileEthereumWallet, referrerEthereumWallet, newProfileGamerTag, &isAfterCreateCalled, txResult, err)
		})
	}
}
