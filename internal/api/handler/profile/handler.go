package profile

import (
	"github.com/kyamalabs/users/api/pb"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/util"
)

type Handler struct {
	pb.UnimplementedProfilesServer
	config util.Config
	store  db.Store
}

func NewHandler(config util.Config, store db.Store) Handler {
	return Handler{
		config: config,
		store:  store,
	}
}
