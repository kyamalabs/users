package db

import "context"

type CreateProfileTxParams struct {
	CreateProfileParams
	AfterCreate func() error
}

type CreateProfileTxResult struct {
	Profile Profile
}

func (store *SQLStore) CreateProfileTx(ctx context.Context, params CreateProfileTxParams) (CreateProfileTxResult, error) {
	var result CreateProfileTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Profile, err = q.CreateProfile(ctx, params.CreateProfileParams)
		if err != nil {
			return err
		}

		return params.AfterCreate()
	})

	return result, err
}
