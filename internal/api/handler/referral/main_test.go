package referral

import (
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/util"
)

func newTestHandler(store db.Store) Handler {
	config := util.Config{}

	return NewHandler(config, store)
}
