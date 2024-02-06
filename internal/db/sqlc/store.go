package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	GetProfileTx(ctx context.Context, params GetProfileTxParams) (GetProfileTxResult, error)
	UpdateProfileTx(ctx context.Context, params UpdateProfileTxParams) (UpdateProfileTxResult, error)
	CreateProfileTx(ctx context.Context, params CreateProfileTxParams) (CreateProfileTxResult, error)
}

type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
