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

const (
	defaultPageSize int32 = 30
)

func (h *Handler) ListReferrals(ctx context.Context, req *pb.ListReferralsRequest) (*pb.ListReferralsResponse, error) {
	logger := log.With().Str("wallet_address", req.GetWalletAddress()).Logger()

	violations := validateGetReferralsRequest(req)
	if violations != nil {
		return nil, handler.InvalidArgumentError(violations)
	}

	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	limit := req.GetPageSize()
	if limit <= 0 {
		limit = defaultPageSize
	}

	offset := (page - 1) * limit

	params := db.ListReferralsParams{
		Limit:    limit,
		Offset:   offset,
		Referrer: req.GetWalletAddress(),
	}

	referrals, err := h.store.ListReferrals(ctx, params)
	if err != nil {
		logger.Error().Err(err).Msg("could not list user referrals")
		return nil, status.Error(codes.Internal, handler.InternalServerError)
	}

	totalReferralsCount, err := h.store.GetReferralsCount(ctx, req.GetWalletAddress())
	if err != nil {
		log.Error().Err(err).Msg("could not get total referrals count")
		return nil, status.Error(codes.Internal, handler.InternalServerError)
	}

	var pbReferrals []*pb.Referral
	for _, referral := range referrals {
		pbReferrals = append(pbReferrals, &pb.Referral{
			Referrer:   referral.Referrer,
			Referee:    referral.Referee,
			ReferredAt: timestamppb.New(referral.ReferredAt),
		})
	}

	response := &pb.ListReferralsResponse{
		Page:           page,
		PageSize:       limit,
		TotalReferrals: int32(totalReferralsCount),
		Referrals:      pbReferrals,
	}

	return response, nil
}

func validateGetReferralsRequest(req *pb.ListReferralsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWalletAddress(req.GetWalletAddress()); err != nil {
		violations = append(violations, handler.FieldViolation("wallet_address", err))
	}

	if err := validator.ValidatePageSize(req.GetPageSize()); err != nil {
		violations = append(violations, handler.FieldViolation("page_size", err))
	}

	return violations
}
