package profile

import (
	"context"

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

func (h *Handler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	logger := log.With().Str("wallet_address", req.GetWalletAddress()).Logger()

	violations := validateUpdateProfileRequest(req)
	if violations != nil {
		return nil, handler.InvalidArgumentError(violations)
	}

	_, err := middleware.AuthorizeUser(ctx, req.GetWalletAddress(), h.authService)
	if err != nil {
		logger.Error().Err(err).Msg("could not authorize user")
		return nil, status.Error(codes.Unauthenticated, handler.UnauthorizedAccessError)
	}

	params := db.UpdateProfileTxParams{
		WalletAddress: req.GetWalletAddress(),
		GamerTag:      req.GetGamerTag(),
		AfterCreate: func() (string, error) {
			return getCachedENSName(ctx, req.GetWalletAddress(), h.cache, h.taskDistributor)
		},
	}

	txResult, err := h.store.UpdateProfileTx(ctx, params)
	if err != nil {
		logger.Error().Err(err).Msg("could not update user profile")
		return nil, status.Error(codes.Internal, handler.InternalServerError)
	}

	response := &pb.UpdateProfileResponse{
		Profile: &pb.Profile{
			WalletAddress: txResult.Profile.WalletAddress,
			GamerTag:      txResult.Profile.GamerTag,
			EnsName:       txResult.EnsName,
			CreatedAt:     timestamppb.New(txResult.Profile.CreatedAt),
		},
	}

	logger.Info().Msg("user profile updated successfully")

	return response, nil
}

func validateUpdateProfileRequest(req *pb.UpdateProfileRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWalletAddress(req.GetWalletAddress()); err != nil {
		violations = append(violations, handler.FieldViolation("wallet_address", err))
	}

	if err := validator.ValidateGamerTag(req.GetGamerTag()); err != nil {
		violations = append(violations, handler.FieldViolation("gamer_tag", err))
	}

	return violations
}
