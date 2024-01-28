package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestGrpcExtractMetadata(t *testing.T) {
	testCases := []struct {
		name        string
		ctx         context.Context
		ctxKey      ReqContextKey
		expectedCtx context.Context
	}{
		{
			name:        "set x-forwarded-for header",
			ctx:         metadata.NewIncomingContext(context.Background(), metadata.Pairs(xForwardedForHeader, "testClientIP")),
			ctxKey:      ClientIP,
			expectedCtx: context.WithValue(context.Background(), ClientIP, "testClientIP"),
		},
		{
			name:        "set x-service-authentication header",
			ctx:         metadata.NewIncomingContext(context.Background(), metadata.Pairs(xServiceAuthenticationHeader, "some service auth header")),
			ctxKey:      ServiceAuthentication,
			expectedCtx: context.WithValue(context.Background(), ServiceAuthentication, "some service auth header"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GrpcExtractMetadata(tc.ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
				require.Equal(t, tc.expectedCtx.Value(tc.ctxKey), ctx.Value(tc.ctxKey))

				return nil, nil
			})

			require.NoError(t, err)
		})
	}
}

func TestHTTPExtractMetadata(t *testing.T) {
	tests := []struct {
		name                          string
		headers                       map[string]string
		expectedClientIP              string
		expectedServiceAuthentication string
	}{
		{
			name: "set x-forwarded-for header",
			headers: map[string]string{
				xForwardedForHeader: "testClientIP",
			},
			expectedClientIP: "testClientIP",
		},
		{
			name: "set x-service-authentication header",
			headers: map[string]string{
				xServiceAuthenticationHeader: "some service auth header",
			},
			expectedServiceAuthentication: "some service auth header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)

			for key, val := range tt.headers {
				req.Header.Set(key, val)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				clientIP := r.Context().Value(ClientIP)
				serviceAuthentication := r.Context().Value(ServiceAuthentication)

				if tt.expectedClientIP != "" {
					require.Equal(t, tt.expectedClientIP, clientIP)
				}

				if tt.expectedServiceAuthentication != "" {
					require.Equal(t, tt.expectedServiceAuthentication, serviceAuthentication)
				}
			})

			httpHandler := HTTPExtractMetadata(handler)
			httpHandler.ServeHTTP(rr, req)
		})
	}
}
