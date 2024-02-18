package referral

import (
	"context"

	"github.com/kyamalabs/users/api/pb"
	"github.com/kyamalabs/users/internal/api/handler"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/validator"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) GetReferrer(ctx context.Context, req *pb.GetReferrerRequest) (*pb.GetReferrerResponse, error) {
	logger := log.With().Str("wallet_address", req.GetWalletAddress()).Logger()

	violations := validateGetReferrerRequest(req)
	if violations != nil {
		return nil, handler.InvalidArgumentError(violations)
	}

	referral, err := h.store.GetReferrer(ctx, req.GetWalletAddress())
	if err != nil && err != db.RecordNotFoundError {
		logger.Error().Err(err).Msg("could not get user referrer")
		return nil, status.Error(codes.Internal, handler.InternalServerError)
	}

	response := &pb.GetReferrerResponse{
		Referral: &pb.Referral{
			Referrer:   referral.Referrer,
			Referee:    referral.Referee,
			ReferredAt: timestamppb.New(referral.ReferredAt),
		},
	}

	return response, nil
}

func validateGetReferrerRequest(req *pb.GetReferrerRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWalletAddress(req.GetWalletAddress()); err != nil {
		violations = append(violations, handler.FieldViolation("wallet_address", err))
	}

	return violations
}
