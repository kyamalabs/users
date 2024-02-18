package profile

import (
	"github.com/kyamalabs/users/internal/cache"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/services"
	"github.com/kyamalabs/users/internal/util"
	"github.com/kyamalabs/users/internal/worker"
)

func newTestHandler(store db.Store, cache cache.Cache, authService services.AuthGrpcService, taskDistributor worker.TaskDistributor) Handler {
	config := util.Config{}

	return NewHandler(config, cache, store, taskDistributor, authService)
}
