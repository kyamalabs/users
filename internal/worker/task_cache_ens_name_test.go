package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	mockcache "github.com/kyamalabs/users/internal/cache/mock"
	"go.uber.org/mock/gomock"

	"github.com/kyamalabs/users/internal/cache"
	"github.com/stretchr/testify/require"
)

func TestTaskCacheENSName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long running async worker test")
	}

	testCases := []struct {
		name            string
		walletAddress   string
		expectedENSName string
	}{
		{
			name:            "success",
			walletAddress:   "0xb865c6093aeAd3c557C32b27033e6C048D2DeB58",
			expectedENSName: "mamabear.eth",
		},
		{
			name:            "ens name not registered for wallet address",
			walletAddress:   "0xA0A17a956A0F59Bd31e6CBEa9d3b8F5602eaD2d2",
			expectedENSName: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testRedisCache.Del(context.Background(), getCacheKey(tc.walletAddress))
			if err != nil && err != cache.Nil {
				require.Fail(t, err.Error())
			}

			payload := &PayloadCacheEnsName{
				WalletAddress: tc.walletAddress,
			}
			err = testTaskDistributor.DistributeTaskCacheEnsName(context.Background(), payload)
			require.NoError(t, err)

			time.Sleep(3 * time.Second) // wait for the task to be processed

			cachedENSName, err := GetCachedENSName(context.Background(), testRedisCache, tc.walletAddress)
			require.NoError(t, err)

			require.Equal(t, tc.expectedENSName, cachedENSName)
		})
	}
}

func TestGetCachedENSName(t *testing.T) {
	testCases := []struct {
		name            string
		walletAddress   string
		buildStubs      func(cache *mockcache.MockCache)
		expectedENSName string
		expectedError   error
		expectedToError bool
	}{
		{
			name:          "success",
			walletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
			buildStubs: func(cache *mockcache.MockCache) {
				cache.EXPECT().
					Get(gomock.Any(), getCacheKey("0xc0ffee254729296a45a3885639AC7E10F9d54979")).
					Times(1).
					Return("bulba", nil)
			},
			expectedENSName: "bulba",
			expectedToError: false,
		},
		{
			name:          "ens name not cached",
			walletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
			buildStubs: func(cache *mockcache.MockCache) {
				cache.EXPECT().
					Get(gomock.Any(), getCacheKey("0xc0ffee254729296a45a3885639AC7E10F9d54979")).
					Times(1).
					Return(nil, nil)
			},
			expectedENSName: "",
			expectedError:   cache.Nil,
			expectedToError: true,
		},
		{
			name:          "error fetching from cache",
			walletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
			buildStubs: func(cache *mockcache.MockCache) {
				cache.EXPECT().
					Get(gomock.Any(), getCacheKey("0xc0ffee254729296a45a3885639AC7E10F9d54979")).
					Times(1).
					Return(nil, errors.New("some cache error"))
			},
			expectedENSName: "",
			expectedToError: true,
		},
		{
			name:          "could not cast cached ens name to string",
			walletAddress: "0xc0ffee254729296a45a3885639AC7E10F9d54979",
			buildStubs: func(cache *mockcache.MockCache) {
				cache.EXPECT().
					Get(gomock.Any(), getCacheKey("0xc0ffee254729296a45a3885639AC7E10F9d54979")).
					Times(1).
					Return(struct{}{}, nil)
			},
			expectedENSName: "",
			expectedToError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := mockcache.NewMockCache(ctrl)

			tc.buildStubs(c)

			cachedENSName, err := GetCachedENSName(context.Background(), c, tc.walletAddress)
			if tc.expectedError != nil {
				require.Equal(t, tc.expectedError, err)
			}
			if tc.expectedToError {
				require.Error(t, err)
			}
			require.Equal(t, tc.expectedENSName, cachedENSName)
		})
	}
}
