package server

import (
	"fmt"
	"sync"

	"github.com/kyamalabs/users/internal/api/handler/referral"

	"github.com/kyamalabs/users/internal/cache"

	"github.com/kyamalabs/users/internal/worker"

	"github.com/kyamalabs/users/internal/api/middleware"
	"github.com/ulule/limiter/v3"

	"github.com/kyamalabs/users/internal/api/handler/profile"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/services"
	"github.com/kyamalabs/users/internal/util"
)

type Server struct {
	ProfileHandler  profile.Handler
	ReferralHandler referral.Handler
}

var once sync.Once

func NewServer(config util.Config, cache cache.Cache, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	authService, err := services.NewAuthServiceGrpcClient(config.AuthServiceGRPCServerAddress, config.ServiceAuthPrivateKeys)
	if err != nil {
		return nil, fmt.Errorf("could not create auth service gRPC client: %w", err)
	}

	err = setupRateLimiter(config.RedisConnURL)
	if err != nil {
		return nil, err
	}

	server := &Server{
		ProfileHandler:  profile.NewHandler(config, cache, store, taskDistributor, authService),
		ReferralHandler: referral.NewHandler(config, store),
	}

	return server, nil
}

func setupRateLimiter(redisConnURL string) error {
	var store limiter.Store
	var createLimiterRedisStoreErr, initializeLimitersErr error

	once.Do(func() {
		store, createLimiterRedisStoreErr = middleware.CreateLimiterRedisStore(redisConnURL)
		if createLimiterRedisStoreErr == nil {
			initializeLimitersErr = middleware.InitializeLimiters(store)
		}
	})

	if createLimiterRedisStoreErr != nil {
		return fmt.Errorf("could not create limiter redis client: %w", createLimiterRedisStoreErr)
	}
	if initializeLimitersErr != nil {
		return fmt.Errorf("could not initialize rate limiters: %w", initializeLimitersErr)
	}

	return nil
}
