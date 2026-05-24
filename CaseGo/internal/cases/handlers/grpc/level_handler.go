package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
	pb "github.com/sewaustav/CaseGogRPServer/gen/go/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LevelGrpcClient interface {
	LevelDoneGrpcHandler(ctx context.Context, user models.UserIdentity, req *dto.LevelXpResult) (*dto.UserLevelInfo, error)
}

type LevelGrpcHandler struct {
	client pb.LevelClient
}

func NewLevelGrpcHandler(addr string) (*LevelGrpcHandler, error) {
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		return nil, fmt.Errorf("read ca cert: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append ca cert")
	}

	// Client certificate
	clientCert, err := tls.LoadX509KeyPair(
		"general-client.crt",
		"general-client.key",
	)
	if err != nil {
		return nil, fmt.Errorf("load client cert: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caPool,

		ServerName: "level-service",

		MinVersion: tls.VersionTLS12,
	}

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(
			credentials.NewTLS(tlsConfig),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}

	client := pb.NewLevelClient(conn)

	return &LevelGrpcHandler{
		client: client,
	}, nil
}

func (h *LevelGrpcHandler) LevelDoneGrpcHandler(ctx context.Context, user models.UserIdentity, req *dto.LevelXpResult) (*dto.UserLevelInfo, error) {
	userProfile, err := h.client.LevelDone(ctx, &pb.LevelResult{
		UserId:     user.UserID,
		Xp:         req.Xp,
		DatePassed: timestamppb.New(time.Now()),
	})
	if err != nil {
		return nil, err
	}
	return &dto.UserLevelInfo{
		Level:          userProfile.Level,
		Xp:             userProfile.Xp,
		Streak:         userProfile.Streak,
		LastActiveDate: userProfile.LastActive.AsTime(),
		IsLevelUp:      userProfile.IsLevelUp,
	}, nil
}
