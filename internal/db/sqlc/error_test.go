package db

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestErrorCode(t *testing.T) {
	testCases := []struct {
		name          string
		buildPgErr    func() error
		expectedError *Error
	}{
		{
			name: "with SQLSTATE code",
			buildPgErr: func() error {
				return &pgconn.PgError{
					Code:           "23505",
					ConstraintName: "table_pkey",
				}
			},
			expectedError: &Error{
				Code:           UniqueViolationCode,
				ConstraintName: "table_pkey",
			},
		},
		{
			name: "without SQLSTATE code",
			buildPgErr: func() error {
				return errors.New("some random error")
			},
			expectedError: &Error{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pgErr := tc.buildPgErr()
			dbError := ParseError(pgErr)

			require.Equal(t, tc.expectedError, dbError)
		})
	}
}
