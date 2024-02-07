package profile

import (
	"github.com/kyamalabs/users/api/pb"
	"github.com/kyamalabs/users/internal/cache"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/services"
	"github.com/kyamalabs/users/internal/util"
	"github.com/kyamalabs/users/internal/worker"
)

type Handler struct {
	pb.UnimplementedProfilesServer
	config          util.Config
	cache           cache.Cache
	store           db.Store
	taskDistributor worker.TaskDistributor
	authService     services.AuthGrpcService
}

func NewHandler(config util.Config, cache cache.Cache, store db.Store, taskDistributor worker.TaskDistributor, authService services.AuthGrpcService) Handler {
	return Handler{
		config:          config,
		cache:           cache,
		store:           store,
		taskDistributor: taskDistributor,
		authService:     authService,
	}
}
