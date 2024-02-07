package middleware

import (
	"context"
	"net/http"

	"github.com/kyamalabs/users/internal/constants"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GrpcExtractMetadata(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		clientIPs := md.Get(xForwardedForHeader)
		if len(clientIPs) > 0 {
			ctx = context.WithValue(ctx, ClientIP, clientIPs[0])
		}

		serviceAuthentications := md.Get(constants.XServiceAuthenticationHeader)
		if len(serviceAuthentications) > 0 {
			ctx = context.WithValue(ctx, ServiceAuthentication, serviceAuthentications[0])
		}
	}

	result, err := handler(ctx, req)

	return result, err
}

func HTTPExtractMetadata(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if xForwardedForHeaderVal := req.Header.Get(xForwardedForHeader); xForwardedForHeaderVal != "" {
			req = req.WithContext(context.WithValue(req.Context(), ClientIP, xForwardedForHeaderVal))
		}

		if xServiceAuthenticationHeaderVal := req.Header.Get(constants.XServiceAuthenticationHeader); xServiceAuthenticationHeaderVal != "" {
			req = req.WithContext(context.WithValue(req.Context(), ServiceAuthentication, xServiceAuthenticationHeaderVal))
		}

		handler.ServeHTTP(res, req)
	})
}
