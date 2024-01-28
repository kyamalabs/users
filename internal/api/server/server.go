package server

import (
	"github.com/kyamalabs/users/internal/api/handler/profile"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/util"
)

type Server struct {
	ProfileHandler profile.Handler
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		ProfileHandler: profile.NewHandler(config, store),
	}

	return server, nil
}
