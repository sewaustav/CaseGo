package rs256

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GrpcInterceptor — это UnaryInterceptor для gRPC
func (m *JWTAuthMiddleware) GrpcInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	authHeader := values[0]
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, status.Error(codes.Unauthenticated, "invalid auth header format")
	}

	claims, err := m.verifyToken(parts[1])
	if err != nil {
		m.logger.Info("token verification failed", "err", err)
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	newCtx := context.WithValue(ctx, UserIDKey, claims.UserID)
	newCtx = context.WithValue(newCtx, RoleKey, claims.Role)

	// 5. Идем дальше
	return handler(newCtx, req)
}
