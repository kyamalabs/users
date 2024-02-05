package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const UniqueViolation = "23505"

var RecordNotFoundError = pgx.ErrNoRows

// ErrorCode returns the PostgreSQL “SQLSTATE” code for a given error if exists; otherwise "".
// see: https://www.postgresql.org/docs/11/errcodes-appendix.html
func ErrorCode(err error) string {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		return pgErr.Code
	}

	return ""
}
