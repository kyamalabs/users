package db

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestErrorCode(t *testing.T) {
	testCases := []struct {
		name              string
		buildPgErr        func() error
		expectedErrorCode string
	}{
		{
			name: "with SQLSTATE code",
			buildPgErr: func() error {
				return &pgconn.PgError{
					Code: "23505",
				}
			},
			expectedErrorCode: UniqueViolation,
		},
		{
			name: "without SQLSTATE code",
			buildPgErr: func() error {
				return errors.New("some random error")
			},
			expectedErrorCode: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pgErr := tc.buildPgErr()
			errorCode := ErrorCode(pgErr)

			require.Equal(t, tc.expectedErrorCode, errorCode)
		})
	}
}
