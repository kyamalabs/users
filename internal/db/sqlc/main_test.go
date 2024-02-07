package db

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kyamalabs/users/internal/util"
	"github.com/rs/zerolog/log"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../../")
	if err != nil {
		log.Fatal().Err(err).Msg("could not load config")
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to the db")
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
