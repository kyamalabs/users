package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateProfileTxParams struct {
	WalletAddress string
	GamerTag      string
	AfterCreate   func() (string, error)
}

type UpdateProfileTxResult struct {
	EnsName string
	Profile Profile
}

func (store *SQLStore) UpdateProfileTx(ctx context.Context, params UpdateProfileTxParams) (UpdateProfileTxResult, error) {
	var result UpdateProfileTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		arg := UpdateProfileParams{
			WalletAddress: params.WalletAddress,
			GamerTag: pgtype.Text{
				String: params.GamerTag,
				Valid:  params.GamerTag != "",
			},
		}

		result.Profile, err = q.UpdateProfile(ctx, arg)
		if err != nil {
			return err
		}

		result.EnsName, err = params.AfterCreate()

		return err
	})

	return result, err
}
