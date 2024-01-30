package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/kyamalabs/users/internal/constants"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kyamalabs/users/api/pb"
	"github.com/kyamalabs/users/internal/api/server"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/kyamalabs/users/docs/statik"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	setupLogger(config)

	runDBMigration(config.DBMigrationURL, config.DBSource)

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	store := db.NewStore(connPool)

	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func setupLogger(config util.Config) {
	logger := log.Logger

	if config.Environment == "development" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger = logger.With().Str("service", constants.ServiceName).Logger()
	log.Logger = logger
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migration instance")
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runGrpcServer(config util.Config, store db.Store) {
	s, err := server.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	grpcInterceptor := grpc.ChainUnaryInterceptor()

	grpcServer := grpc.NewServer(grpcInterceptor)
	pb.RegisterProfilesServer(grpcServer, s.ProfileHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create grpc server listener")
	}

	log.Info().Msgf("started gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	s, err := server.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	opt := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(opt)

	ctx, cancel := context.WithCancel(context.Background())

	err = pb.RegisterProfilesHandlerServer(ctx, grpcMux, s.ProfileHandler)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register profiles handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	srv := &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create http gateway server listener")
	}

	log.Info().Msgf("started HTTP gateway server at %s", listener.Addr().String())

	err = srv.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server")
	}

	cancel()
}
