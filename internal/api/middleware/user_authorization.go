package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kyamalabs/users/internal/constants"

	"github.com/kyamalabs/auth/api/pb"
	"github.com/kyamalabs/users/internal/services"
	"google.golang.org/grpc/metadata"
)

func AuthorizeUser(ctx context.Context, walletAddress string, authService services.AuthGrpcService) (*pb.AccessTokenPayload, error) {
	mtdt, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("could not get metadata from incoming context")
	}

	authValues := mtdt.Get(constants.AuthorizationHeader)
	if len(authValues) == 0 {
		return nil, errors.New("missing authorization header")
	}

	authHeader := authValues[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, errors.New("invalid authorization header format")
	}

	authType := fields[0]
	if !strings.EqualFold(constants.AuthorizationBearer, authType) {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	payload := &pb.VerifyAccessTokenRequest{
		WalletAddress: walletAddress,
	}
	response, err := authService.VerifyAccessToken(payload, accessToken)
	if err != nil {
		return nil, fmt.Errorf("could not verify access token: %w", err)
	}

	return response.GetPayload(), nil
}
