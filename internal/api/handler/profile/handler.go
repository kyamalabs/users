package profile

import (
	"github.com/kyamalabs/users/api/pb"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/services"
	"github.com/kyamalabs/users/internal/util"
)

type Handler struct {
	pb.UnimplementedProfilesServer
	config      util.Config
	store       db.Store
	authService services.AuthGrpcService
}

func NewHandler(config util.Config, store db.Store, authService services.AuthGrpcService) Handler {
	return Handler{
		config:      config,
		store:       store,
		authService: authService,
	}
}
