package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateWalletAddress(t *testing.T) {
	testCases := []struct {
		name            string
		walletAddress   string
		expectedToError bool
	}{
		{
			name:            "no error",
			walletAddress:   "0x999999cf1046e68e36E1aA2E0E07105eDDD1f08E",
			expectedToError: false,
		},
		{
			name:            "not a valid hex address",
			walletAddress:   "0x999999cf1046e68e36E1aA2E0E07105eDDD1f08",
			expectedToError: true,
		},
		{
			name:            "not prefixed with 0x",
			walletAddress:   "999999cf1046e68e36E1aA2E0E07105eDDD1f08E",
			expectedToError: true,
		},
		{
			name:            "is a zero address",
			walletAddress:   "0x0000000000000000000000000000000000000000",
			expectedToError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateWalletAddress(tc.walletAddress)
			if tc.expectedToError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestValidateGamerTag(t *testing.T) {
	testCases := []struct {
		name            string
		gamerTag        string
		expectedToError bool
	}{
		{
			name:            "no error",
			gamerTag:        "us3r_mx_tt43",
			expectedToError: false,
		},
		{
			name:            "too long",
			gamerTag:        "thi3_gam3r_5ag_is_t00_l0n6",
			expectedToError: true,
		},
		{
			name:            "too short",
			gamerTag:        "69",
			expectedToError: true,
		},
		{
			name:            "invalid characters",
			gamerTag:        "@inval#d",
			expectedToError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateGamerTag(tc.gamerTag)
			if tc.expectedToError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
