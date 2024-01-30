package server

import (
	"fmt"

	"github.com/kyamalabs/users/internal/api/handler/profile"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/services"
	"github.com/kyamalabs/users/internal/util"
)

type Server struct {
	ProfileHandler profile.Handler
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	authService, err := services.NewAuthServiceGrpcClient(config.AuthServiceGRPCServerAddress, config.ServiceAuthPrivateKeys)
	if err != nil {
		return nil, fmt.Errorf("could not create auth service gRPC client: %w", err)
	}

	server := &Server{
		ProfileHandler: profile.NewHandler(config, store, authService),
	}

	return server, nil
}
