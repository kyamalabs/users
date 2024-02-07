package profile

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

func (h *Handler) ListAllProfiles(ctx context.Context, req *pb.ListAllProfilesRequest) (*pb.ListAllProfilesResponse, error) {
	violations := validateListAllProfileRequest(req)
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

	params := db.ListProfilesParams{
		Limit:  limit,
		Offset: offset,
	}

	profiles, err := h.store.ListProfiles(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("could not list user profiles")
		return nil, status.Error(codes.Internal, handler.InternalServerError)
	}

	totalProfilesCount, err := h.store.GetProfilesCount(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not get total profiles count")
		return nil, status.Error(codes.Internal, handler.InternalServerError)
	}

	var publicProfiles []*pb.PublicProfile
	for _, profile := range profiles {
		ensName, err := getCachedENSName(ctx, profile.WalletAddress, h.cache, h.taskDistributor)
		if err != nil {
			log.Error().Err(err).Msg("could not get cached ens name")
		}

		publicProfiles = append(publicProfiles, &pb.PublicProfile{
			WalletAddress: profile.WalletAddress,
			GamerTag:      profile.GamerTag,
			EnsName:       ensName,
			CreatedAt:     timestamppb.New(profile.CreatedAt),
		})
	}

	response := &pb.ListAllProfilesResponse{
		Page:          page,
		PageSize:      limit,
		TotalProfiles: int32(totalProfilesCount),
		Profiles:      publicProfiles,
	}

	return response, nil
}

func validateListAllProfileRequest(req *pb.ListAllProfilesRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePageSize(req.GetPageSize()); err != nil {
		violations = append(violations, handler.FieldViolation("page_size", err))
	}

	return violations
}
