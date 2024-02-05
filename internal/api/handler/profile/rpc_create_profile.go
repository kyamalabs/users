package profile

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/kyamalabs/users/internal/worker"

	"github.com/kyamalabs/users/api/pb"
	"github.com/kyamalabs/users/internal/api/handler"
	"github.com/kyamalabs/users/internal/api/middleware"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/validator"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	logger := log.With().Str("wallet_address", req.GetWalletAddress()).Logger()

	violations := validateCreateProfileRequest(req)
	if violations != nil {
		return nil, handler.InvalidArgumentError(violations)
	}

	_, err := middleware.AuthorizeUser(ctx, req.GetWalletAddress(), h.authService)
	if err != nil {
		logger.Error().Err(err).Msg("could not authorize user")
		return nil, status.Error(codes.Unauthenticated, handler.UnauthorizedAccessError)
	}

	params := db.CreateProfileTxParams{
		CreateProfileParams: db.CreateProfileParams{
			WalletAddress: req.GetWalletAddress(),
			GamerTag:      req.GetGamerTag(),
		},
		AfterCreate: func() error {
			return cacheENSName(ctx, req.GetWalletAddress(), h.taskDistributor)
		},
	}

	txResult, err := h.store.CreateProfileTx(ctx, params)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			logger.Error().Err(err).Msg("user profile already exists")
			return nil, status.Error(codes.AlreadyExists, AlreadyExists)
		}

		logger.Error().Err(err).Msg("could not create user profile")
		return nil, status.Error(codes.Internal, handler.InternalServerError)
	}

	response := &pb.CreateProfileResponse{
		Profile: &pb.Profile{
			WalletAddress: txResult.Profile.WalletAddress,
			GamerTag:      txResult.Profile.GamerTag,
			CreatedAt:     timestamppb.New(txResult.Profile.CreatedAt),
		},
	}

	logger.Info().Msg("user profile created successfully")

	return response, nil
}

func cacheENSName(ctx context.Context, walletAddress string, taskDistributor worker.TaskDistributor) error {
	taskPayload := &worker.PayloadCacheEnsName{
		WalletAddress: walletAddress,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	return taskDistributor.DistributeTaskCacheEnsName(ctx, taskPayload, opts...)
}

func validateCreateProfileRequest(req *pb.CreateProfileRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWalletAddress(req.GetWalletAddress()); err != nil {
		violations = append(violations, handler.FieldViolation("wallet_address", err))
	}

	if err := validator.ValidateGamerTag(req.GetGamerTag()); err != nil {
		violations = append(violations, handler.FieldViolation("gamer_tag", err))
	}

	return violations
}
