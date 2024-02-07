package profile

import (
	"context"

	"github.com/kyamalabs/users/api/pb"
	"github.com/kyamalabs/users/internal/api/handler"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) GetPublicProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetPublicProfileResponse, error) {
	logger := log.With().Str("wallet_address", req.GetWalletAddress()).Logger()

	violations := validateGetProfileRequest(req)
	if violations != nil {
		return nil, handler.InvalidArgumentError(violations)
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

	response := &pb.GetPublicProfileResponse{
		Profile: &pb.PublicProfile{
			WalletAddress: txResult.Profile.WalletAddress,
			EnsName:       txResult.EnsName,
			GamerTag:      txResult.Profile.GamerTag,
			CreatedAt:     timestamppb.New(txResult.Profile.CreatedAt),
		},
	}

	logger.Info().Msg("fetched user profile successfully")

	return response, nil
}
