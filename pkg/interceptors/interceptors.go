package interceptors

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

const (
	requestUUIDKey = "requestUUID"
)

func UnaryRequestUUIDGenerator() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (
		interface{}, error,
	) {
		requestUUID := uuid.NewV4()
		ctx = context.WithValue(ctx, requestUUIDKey, requestUUID.String())

		resp, err := handler(ctx, req)

		return resp, err
	}
}
