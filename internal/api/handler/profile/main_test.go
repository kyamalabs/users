package profile

import (
	"context"
	"fmt"
	"testing"

	"github.com/kyamalabs/users/internal/cache"
	db "github.com/kyamalabs/users/internal/db/sqlc"
	"github.com/kyamalabs/users/internal/services"
	"github.com/kyamalabs/users/internal/util"
	"github.com/kyamalabs/users/internal/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func newTestHandler(store db.Store, cache cache.Cache, authService services.AuthGrpcService, taskDistributor worker.TaskDistributor) Handler {
	config := util.Config{}

	return NewHandler(config, cache, store, taskDistributor, authService)
}

func newContextWithBearerToken() context.Context {
	bearerToken := fmt.Sprintf("%s %s", "bearer", "some-token")
	md := metadata.MD{
		"Authorization": []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}

func checkInvalidRequestParams(t *testing.T, err error, expectedFieldViolations []string) {
	var violations []string

	st, ok := status.FromError(err)
	require.True(t, ok)

	details := st.Details()

	for _, detail := range details {
		br, ok := detail.(*errdetails.BadRequest)
		require.True(t, ok)

		fieldViolations := br.FieldViolations
		for _, violation := range fieldViolations {
			violations = append(violations, violation.Field)
		}
	}

	require.ElementsMatch(t, expectedFieldViolations, violations)
}
