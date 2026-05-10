package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	pb "github.com/sewaustav/CaseGogRPServer/gen/go/payments"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type PaymentGrpcClient interface {
	CheckStatus(ctx context.Context, userID int64) (*dto.SubscriptionStatusDto, error)
}

type PaymentCheckGRPCHandler struct {
	client pb.PaymentCheckServiceClient
}

func NewPaymentCheckGRPCHandler(addr string) (PaymentGrpcClient, error) {
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

		ServerName: "payment-service",

		MinVersion: tls.VersionTLS12,
	}

	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(
			credentials.NewTLS(tlsConfig),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}

	client := pb.NewPaymentCheckServiceClient(conn)

	return &PaymentCheckGRPCHandler{
		client: client,
	}, nil
}

func (h *PaymentCheckGRPCHandler) CheckStatus(ctx context.Context, userID int64) (*dto.SubscriptionStatusDto, error) {
	status, err := h.client.CheckStatus(ctx, &pb.UserInfo{
		UserId: userID,
	})
	if err != nil {
		return nil , err 
	}

	return &dto.SubscriptionStatusDto{
		Status: status.Status,
		ExpiredAt: status.ExpiredAt.AsTime(),
	}, nil
}