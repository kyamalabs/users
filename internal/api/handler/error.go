package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	InternalServerError     string = "An unexpected error occurred while processing your request."
	UnauthorizedAccessError string = "Authorization failed. Please verify your credentials and try again."
)

func FieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func InvalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")

	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return statusDetails.Err()
}

func CheckInvalidRequestParams(t *testing.T, err error, expectedFieldViolations []string) {
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

func NewContextWithBearerToken() context.Context {
	bearerToken := fmt.Sprintf("%s %s", "bearer", "some-token")
	md := metadata.MD{
		"Authorization": []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}
