package profile

import (
	"context"

	"github.com/kyamalabs/users/api/pb"
	"github.com/kyamalabs/users/internal/api/handler"
	"github.com/kyamalabs/users/internal/api/middleware"
	"github.com/kyamalabs/users/internal/cache"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/validator"
	"github.com/kyamalabs/users/internal/worker"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	logger := log.With().Str("wallet_address", req.GetWalletAddress()).Logger()

	violations := validateGetProfileRequest(req)
	if violations != nil {
		return nil, handler.InvalidArgumentError(violations)
	}

	_, err := middleware.AuthorizeUser(ctx, req.GetWalletAddress(), h.authService)
	if err != nil {
		logger.Error().Err(err).Msg("could not authorize user")
		return nil, status.Error(codes.Unauthenticated, handler.UnauthorizedAccessError)
	}

	params := db.GetProfileTxParams{
		WalletAddress: req.GetWalletAddress(),
		AfterCreate: func() (string, error) {
			return getCachedENSName(ctx, req.GetWalletAddress(), h.cache, h.taskDistributor)
		},
	}

	txResult, err := h.store.GetProfileTx(ctx, params)
	if err != nil {
		if err == db.RecordNotFoundError {
			logger.Error().Err(err).Msg("user profile does not exist")
			return nil, status.Error(codes.NotFound, DoesNotExist)
		}

		logger.Error().Err(err).Msg("could not get user profile")
		return nil, status.Error(codes.Internal, handler.InternalServerError)
	}

	response := &pb.GetProfileResponse{
		Profile: &pb.Profile{
			WalletAddress: txResult.Profile.WalletAddress,
			EnsName:       txResult.EnsName,
			GamerTag:      txResult.Profile.GamerTag,
			CreatedAt:     timestamppb.New(txResult.Profile.CreatedAt),
		},
	}

	logger.Info().Msg("fetched user profile successfully")

	return response, nil
}

func getCachedENSName(ctx context.Context, walletAddress string, c cache.Cache, taskDistributor worker.TaskDistributor) (string, error) {
	ensName, err := worker.GetCachedENSName(ctx, c, walletAddress)
	if err == cache.Nil {
		err = cacheENSName(ctx, walletAddress, taskDistributor)
		return "", err
	}

	return ensName, err
}

func validateGetProfileRequest(req *pb.GetProfileRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWalletAddress(req.GetWalletAddress()); err != nil {
		violations = append(violations, handler.FieldViolation("wallet_address", err))
	}

	return violations
}
