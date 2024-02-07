package profile

import (
	"context"

	"github.com/kyamalabs/users/api/pb"
	"github.com/kyamalabs/users/internal/api/handler"
	"github.com/kyamalabs/users/internal/api/middleware"
	"github.com/kyamalabs/users/internal/validator"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*emptypb.Empty, error) {
	logger := log.With().Str("wallet_address", req.GetWalletAddress()).Logger()

	violations := validateDeleteProfileRequest(req)
	if violations != nil {
		return nil, handler.InvalidArgumentError(violations)
	}

	_, err := middleware.AuthorizeUser(ctx, req.GetWalletAddress(), h.authService)
	if err != nil {
		logger.Error().Err(err).Msg("could not authorize user")
		return nil, status.Error(codes.Unauthenticated, handler.UnauthorizedAccessError)
	}

	err = h.store.DeleteProfile(ctx, req.GetWalletAddress())
	if err != nil {
		logger.Error().Err(err).Msg("could not create user profile")
		return nil, status.Error(codes.Internal, handler.InternalServerError)
	}

	logger.Info().Msg("user profile deleted successfully")

	return &emptypb.Empty{}, nil
}

func validateDeleteProfileRequest(req *pb.DeleteProfileRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWalletAddress(req.GetWalletAddress()); err != nil {
		violations = append(violations, handler.FieldViolation("wallet_address", err))
	}

	return violations
}
