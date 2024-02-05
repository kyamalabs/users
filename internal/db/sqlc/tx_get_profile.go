package db

import "context"

type GetProfileTxParams struct {
	WalletAddress string
	AfterCreate   func() (string, error)
}

type GetProfileTxResult struct {
	EnsName string
	Profile Profile
}

func (store *SQLStore) GetProfileTx(ctx context.Context, params GetProfileTxParams) (GetProfileTxResult, error) {
	var result GetProfileTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Profile, err = q.GetProfile(ctx, params.WalletAddress)
		if err != nil {
			return err
		}

		result.EnsName, err = params.AfterCreate()

		return err
	})

	return result, err
}
